package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type CollectionElem struct {
	Name              string `json:"name"`
	RequiredPeerCount int    `json:"requiredPeerCount"`
	MaxPeerCount      int    `json:"maxPeerCount"`
	BlockToLive       int    `json:"blockToLive"`
	MemberOnlyRead    bool   `json:"memberOnlyRead"`
	Policy            string `json:"policy"`
}

func generateCollection(nOrgs int) {
	collection := []CollectionElem{}

	for _, a := range assetTypeList {
		if len(a.Readers) > 0 {
			elem := CollectionElem{
				Name:              a.Tag,
				RequiredPeerCount: 0,
				MaxPeerCount:      3,
				BlockToLive:       1000000,
				MemberOnlyRead:    true,
				Policy:            generatePolicy(a.Readers, nOrgs),
			}
			collection = append(collection, elem)
		}
	}

	b, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("collections.json", b, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func generatePolicy(readers []string, nOrgs int) string {
	firstElem := true
	policy := "OR("
	for _, r := range readers {
		if nOrgs == 1 {
			if r == "org1MSP" || r == "org2MSP" || r == "org3MSP" {
				continue
			}
		} else if nOrgs == 3 {
			if r == "orgMSP" {
				continue
			}
		}
		if !firstElem {
			policy += ", "
		}
		policy += fmt.Sprintf("'%s.member'", r)
		firstElem = false
	}
	policy += ")"
	return policy
}
