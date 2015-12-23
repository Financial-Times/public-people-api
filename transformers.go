package main

import (
	"github.com/Financial-Times/neoism"
)

// Thing Decorates a result with the base properties of a Thing
func Thing(node *neoism.Node, result *(map[string]interface{})) {
	resMap := *result
	labels, err := node.Labels()
	if err != nil {
		panic(err)
	}
	for key, value := range node.Data {
		resMap[key] = value
	}
	resMap["types"] = labels
	if resMap["factsetIdentifier"] != nil {
		delete(resMap, "factsetIdentifier")
	}
}

// Person decorates the result structure with properties from a Person
func Person(node *neoism.Node, result *map[string]interface{}) {
	resMap := *result
	resMap["plnkyname"] = node.Data["name"]
}
