package core

import (
	"flag"
	"fmt"
	"sync"
	"watch01/assets"
)

var Wg sync.WaitGroup
var URL = flag.String("url", "", "input url")
var Urllist = flag.String("file", "", "input path to urllist.txt")

func Get_pa() {
	flag.Parse()
}
func Run(url string, Bodys string, Headers string, Servers string) {
	jsonConfigList := assets.GetConfigs()
	unmarshelledConfigs := assets.DeserializeJson(jsonConfigList)
	for _, configObj := range unmarshelledConfigs {
		fmt.Println("正在扫描: " + configObj.Product)
		i := 0
		for i = 0; i < len(configObj.Rules); i++ {
			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, Check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, Bodys, Headers, Servers))
				}
				if !(ContainsInSlice(arr, false)) {
					if ContainsInSlice(arr, true) {
						Resulte_write(url, configObj.Product)
					}
				}

			} else {
				if Check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, Bodys, Headers, Servers) {
					Resulte_write(url, configObj.Product)
				}
			}
		}
	}
}
