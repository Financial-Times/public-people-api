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
		R  *neoism.Node
		RR *neoism.Relationship
		mm *neoism.Relationship
		oo *neoism.Relationship
	}{}

	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (p:Thing{uuid: {uuid}})<-[mm:HAS_MEMBER]-(m:Membership)
                        MATCH (m)-[rr:HAS_ROLE]->(r:Role)
                        MATCH (m)-[oo:HAS_ORGANISATION]->(o:Organisation)
                        RETURN m, {o, collect(r)}, mm, rr, oo
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
		result.R.Db = mcd.db
		membership := make(map[string]interface{})
		Thing(result.M, &membership)
		organisation := make(map[string]interface{})
		Thing(result.O, &organisation)
		role := make(map[string]interface{})
		Thing(result.R, &role)
		membership["organisation"] = organisation
		membership["role"] = role
		memberships[idx] = membership
	}
	return memberships, true, nil
}
