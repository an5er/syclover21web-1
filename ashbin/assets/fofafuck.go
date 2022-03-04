package assets

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
)

type Rules [][]struct {
	Match   string `json:"match"`
	Content string `json:"content"`
}

type Fofa_dic struct {
	RuleID         string `json:"rule_id"`
	Level          string `json:"level"`
	Softhard       string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Rules          Rules  `json:"rules"`
}

func GetConfigs() string {
	b, err := ioutil.ReadFile("fofa.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	return str
}

func DeserializeJson(configJson string) []Fofa_dic {
	jsonAsBytes := []byte(configJson)
	configs := make([]Fofa_dic, 0)
	var jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	err := jsons.Unmarshal(jsonAsBytes, &configs)
	if err != nil {
		panic(err)
	}
	return configs
}
