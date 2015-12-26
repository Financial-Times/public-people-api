package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Financial-Times/neoism"
)

// MembershipDriver interface
type MembershipDriver interface {
	FindMembershipsByPersonUUID(uuid string) ([]Membership, error)
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
func (mcd MembershipCypherDriver) FindMembershipsByPersonUUID(uuid string) ([]Membership, error) {

	// result := []struct {
	// 	M struct {
	// 		Membership neoMembership `json:"Data"`
	// 	}
	// }{}

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
                      RETURN m, o, r, mm, rr, oo
                      `,
		Parameters:   neoism.Props{"uuid": uuid},
		Result:       &results,
		IncludeStats: true,
	}

	err := mcd.db.Cypher(query)

	if err != nil {
		return nil, err
	}

	//log.Println(query.Statement)
	//log.Printf("Returned structure %+v\n", results)
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	} else {
		os.Stdout.Write(b)
	}
	for key, result := range results {
		result.M.Db = mcd.db
		labels, err := result.M.Labels()
		properties, err := result.M.Properties()
		relationships, err := result.M.Relationships()
		data := result.O.Data
		if err != nil {
			panic(err)
		}
		//json.NewEncoder(os.Stdout).Encode(results)
		//log.Printf("Returned structure Row %d Membership: %+v, Organisation: %+v, Role: %+v \n", key, result.M, result.O, result.R)
		log.Printf("\nReturned structure Row %d Labels: %+v, Properties: %+v, Relationships %+v, Data: %v \n", key, labels, properties, relationships, data)
	}
	//log.Printf("Labels %+v Data %+v Relationships %+v \n", result[0].N.Labels(), result[0].Data, result[0].HrefAllTypedRels)
	// var membership neoMembership
	var memberships []Membership
	// for key, result := range result {
	// 	membership = result.M.Membership
	// 	memberships[key] = toMembership(membership)
	// }

	return memberships, nil
}
