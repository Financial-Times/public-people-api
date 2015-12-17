package main

import (
	"log"

	"github.com/jmcvetta/neoism"
)

// PeopleDriver interface
type PeopleDriver interface {
	Read(id string) Person
}

// PeopleCypherDriver struct
type PeopleCypherDriver struct {
	db *neoism.Database
}

//NewPeopleCypherDriver instanciate driver
func NewPeopleCypherDriver(db *neoism.Database) PeopleCypherDriver {
	return PeopleCypherDriver{db}
}

func (pcw PeopleCypherDriver) Read(uuid string) Person {

	result := []struct {
		P struct {
			Person neoPerson `json:"Data"`
		}
	}{}

	query := &neoism.CypherQuery{
		Statement: `
      MATCH (p:Person {uuid: {uuid}})
      RETURN p
    `,
		Parameters: neoism.Props{"uuid": uuid},
		Result:     &result,
	}

	err := pcw.db.Cypher(query)

	log.Println(query.Statement)

	if err != nil {
		panic(err)
	}
	log.Printf("Returned structure %+v\n", result[0])
	//log.Printf("Labels %+v Data %+v Relationships %+v \n", result[0].N.Labels(), result[0].Data, result[0].HrefAllTypedRels)
	p := result[0].P.Person
	return toPerson(p)
}
