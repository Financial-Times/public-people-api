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

// NewPeopleCypherDriver instanciate driver
func NewPeopleCypherDriver(db *neoism.Database) PeopleCypherDriver {
	return PeopleCypherDriver{db}
}

func (pcw PeopleCypherDriver) Read(uuid string) Person {
	result := []struct {
		N neoism.Node
	}{}

	query := &neoism.CypherQuery{
		Statement: `
      MATCH (p:Person) WHERE p.uuid = {uuid}
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
	log.Println(result)
	return Person{}

}
