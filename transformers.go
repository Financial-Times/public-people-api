package main

import (
	"fmt"
	"github.com/Financial-Times/neoism"
)

// Thing Decorates a result with the base properties of a Thing
func Thing(node *neoism.Node, result *(map[string]interface{})) {
	resMap := *result
	labels, err := node.Labels()
	if err != nil {
		panic(err)
	}
	resMap["types"] = labels
	for key, value := range node.Data {
		resMap[key] = value
	}
	if resMap["factsetIdentifier"] != nil {
		delete(resMap, "factsetIdentifier")
	} else if resMap["fsIdentifier"] != nil {
		delete(resMap, "fsIdentifier")
	}
	if resMap["uuid"] != nil {
		resMap["uri"] = fmt.Sprintf("http://api.ft.com/things/%s", resMap["uuid"])
		resMap["apiUrl"] = fmt.Sprintf("http://api.ft.com/%s/%s", thingURLType(labels), resMap["uuid"])
		delete(resMap, "uuid")
	}
	changeEvents(result)
}

func changeEvents(res *(map[string]interface{})) {
	resMap := *res
	start := resMap["inceptionDate"]
	end := resMap["terminationDate"]
	if start != nil || end != nil {
		changeEvents := make(map[string]interface{})
		if start != nil {
			changeEvents["started"] = start
			delete(resMap, "inceptionDate")
		}
		if end != nil {
			changeEvents["ended"] = end
			delete(resMap, "terminationDate")
		}
		resMap["changeEvents"] = changeEvents
	}
}

func thingURLType(types []string) string {
	for _, thingType := range types {
		switch thingType {
		case "Person":
			return "people"
		case "Organisation", "Company", "PublicCompany", "PrivateCompany":
			return "organisations"
		case "Membership":
			return "memberships"
		case "Role":
			return "roles"
		}
	}
	return "things"
}
