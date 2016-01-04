package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
)

// Thing builds a result based on the properties of the neoism.Node passed
func Thing(node *neoism.Node, result *(map[string]interface{})) {
	resMap := *result
	labels, err := node.Labels()
	if err != nil {
		panic(err)
	}
	resMap["types"] = typeURIs(labels)
	for key, value := range node.Data {
		resMap[key] = value
	}
	if resMap["uuid"] != nil {
		resMap["id"] = fmt.Sprintf("http://api.ft.com/things/%s", resMap["uuid"])
		resMap["apiUrl"] = fmt.Sprintf("http://api.ft.com/%s/%s", thingURLType(labels), resMap["uuid"])
	}
	changeEvents(result)
	cleanUp(result)
}

func cleanUp(resMap *(map[string]interface{})) {
	delete(*resMap, "uuid")
	delete(*resMap, "factsetIdentifier")
	delete(*resMap, "fsIdentifier")
	delete(*resMap, "inceptionDate")
	delete(*resMap, "terminationDate")
}

func typeURIs(labels []string) []string {
	base := "http://www.ft.com/ontology/"
	var results []string
	for _, label := range labels {
		switch label {
		case "Person":
			results = append(results, base+"person/Person")
			break
		case "Organisation", "Company", "PublicCompany", "PrivateCompany":
			results = append(results, base+"organisation/"+label)
			break
		case "Thing":
			results = append(results, base+"core/Thing")
			results = append(results, base+"core/Concept")
			break
		case "Role":
			results = append(results, base+"organisation/"+label)
			break
		case "Membership":
			results = append(results, base+"organisation/"+label)
			break
		}
	}
	return results
}

func changeEvents(res *(map[string]interface{})) {
	resMap := *res
	start := resMap["inceptionDate"]
	end := resMap["terminationDate"]
	if start != nil || end != nil {
		changeEvents := make(map[string]interface{})
		if start != nil {
			changeEvents["started"] = start
		}
		if end != nil {
			changeEvents["ended"] = end
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
