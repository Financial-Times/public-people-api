package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/jmcvetta/neoism"
)

// PeopleDriver interface
type PeopleDriver interface {
	Read(id string) map[string]interface{}
}

// PeopleCypherDriver struct
type PeopleCypherDriver struct {
	db *neoism.Database
}

//NewPeopleCypherDriver instanciate driver
func NewPeopleCypherDriver(db *neoism.Database) PeopleCypherDriver {
	return PeopleCypherDriver{db}
}

func (pcw PeopleCypherDriver) Read(uuid string) map[string]interface{} {
	results := []struct {
		P *neoism.Node
	}{}
	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (p:Person {uuid: {uuid}})
                        RETURN p
                        `,
		Parameters: neoism.Props{"uuid": uuid},
		Result:     &results,
	}
	err := pcw.db.Cypher(query)
	if err != nil {
		panic(err)
	}
	result := make(map[string]interface{})
	results[0].P.Db = pcw.db
	Thing(results[0].P, &result)
	log.Debugf("Returning %v", result)
	return result
}
