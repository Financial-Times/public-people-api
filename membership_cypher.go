package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// MembershipDriver interface
type MembershipDriver interface {
	FindMembershipsByPersonUUID(uuid string) ([]interface{}, bool, error)
}

// MembershipCypherDriver struct
type MembershipCypherDriver struct {
	db *neoism.Database
}

//NewMembershipCypherDriver instanciate driver
func NewMembershipCypherDriver(db *neoism.Database) MembershipCypherDriver {
	return MembershipCypherDriver{db}
}

//FindMembershipsByPersonUUID returns the memberships associated with a Person
func (mcd MembershipCypherDriver) FindMembershipsByPersonUUID(uuid string) ([]interface{}, bool, error) {

	results := []struct {
		M  *neoism.Node
		O  *neoism.Node
		MM *neoism.Relationship
		OO *neoism.Relationship
		R  []struct {
			R  *neoism.Node
			RR *neoism.Relationship
		}
	}{}

	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (:Thing{uuid: {uuid}})<-[:HAS_MEMBER]-(m:Membership)
                        OPTIONAL MATCH (m)-[:HAS_ORGANISATION]->(o:Organisation)
                        OPTIONAL MATCH (m)-[rr:HAS_ROLE]->(r:Role)
                        RETURN m, o, collect({ r:r, rr:rr}) as r
                        `,
		Parameters: neoism.Props{"uuid": uuid},
		Result:     &results,
	}

	err := mcd.db.Cypher(query)
	if err != nil {
		log.WithField("error", err).Error("Error executing cypher")
		return nil, false, err
	} else if len(results) == 0 {
		log.WithField("uuid", uuid).Debug("No memberships found")
		return nil, false, nil
	}

	log.WithFields(log.Fields{"uuid": uuid, "matches": len(results)}).Info("Memberships found")
	memberships := make([]interface{}, len(results))
	for idx, result := range results {
		result.M.Db = mcd.db
		result.O.Db = mcd.db
		membership := make(map[string]interface{})
		Thing(result.M, &membership)
		organisation := make(map[string]interface{})
		Thing(result.O, &organisation)
		membership["organisation"] = organisation
		roles := make([]interface{}, len(result.R))
		for idx, roleResult := range result.R {
			roleResult.R.Db = mcd.db
			roleResult.RR.Db = mcd.db
			props, _ := roleResult.RR.Properties()
			for key, value := range props {
				roleResult.R.Data[key] = value
			}
			role := make(map[string]interface{})
			Thing(roleResult.R, &role)
			roles[idx] = role
		}
		membership["roles"] = roles
		memberships[idx] = membership
	}
	return memberships, true, nil
}
