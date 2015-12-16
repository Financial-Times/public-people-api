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

// Result struct
// type Result struct {
// 	Columns string `json:"columns"`
// 	Data    string `json:"data"`
// }

//NewPeopleCypherDriver instanciate driver
func NewPeopleCypherDriver(db *neoism.Database) PeopleCypherDriver {
	return PeopleCypherDriver{db}
}

func (pcw PeopleCypherDriver) Read(uuid string) Person {

	result := []struct {
		Person Person `json:"p"`
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
	log.Printf("%+v\n", result)
	return Person{}

}
