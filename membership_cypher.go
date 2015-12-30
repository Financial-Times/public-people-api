package main

import "github.com/Financial-Times/neoism"

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
		mm *neoism.Relationship
		oo *neoism.Relationship
	}{}

	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (p:Thing{uuid: {uuid}})<-[mm:HAS_MEMBER]-(m:Membership)
                        MATCH (m)-[oo:HAS_ORGANISATION]->(o:Organisation)
                        RETURN m, o, mm, oo
                        `,
		Parameters:   neoism.Props{"uuid": uuid},
		Result:       &results,
		IncludeStats: true,
	}

	err := mcd.db.Cypher(query)
	if err != nil {
		return nil, false, err
	} else if len(results) == 0 {
		return nil, false, nil
	}

	memberships := make([]interface{}, len(results))
	for idx, result := range results {
		result.M.Db = mcd.db
		result.O.Db = mcd.db
		membership := make(map[string]interface{})
		Thing(result.M, &membership)
		organisation := make(map[string]interface{})
		Thing(result.O, &organisation)
		membership["organisation"] = organisation
		mUUID, err := result.M.Property("uuid")
		roles, found, err := mcd.findRolesbyMembershipUUID(mUUID)
		if err != nil || !found {
			membership["roles"] = nil
		} else {
			membership["roles"] = roles
		}
		memberships[idx] = membership
	}
	return memberships, true, nil
}

func (mcd MembershipCypherDriver) findRolesbyMembershipUUID(uuid string) ([]interface{}, bool, error) {
	results := []struct {
		R  *neoism.Node
		RR *neoism.Relationship
	}{}

	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (m:Membership{uuid: {uuid}})
                        MATCH (m)-[rr:HAS_ROLE]->(r:Role)
                        RETURN r, rr
                        `,
		Parameters:   neoism.Props{"uuid": uuid},
		Result:       &results,
		IncludeStats: true,
	}

	err := mcd.db.Cypher(query)
	if err != nil {
		return nil, false, err
	} else if len(results) == 0 {
		return nil, false, nil
	}
	roles := make([]interface{}, len(results))
	for idx, result := range results {
		result.R.Db = mcd.db
		role := make(map[string]interface{})
		Thing(result.R, &role)
		roles[idx] = role
	}
	return roles, true, nil
}
