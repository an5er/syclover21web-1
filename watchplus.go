package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var URL = flag.String("url", "", "input url")
var Urllist = flag.String("file", "", "input path to urllist.txt")

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

func getConfigs() string {
	b, err := ioutil.ReadFile("fofa.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	return str
}

func deserializeJson(configJson string) []Fofa_dic {
	jsonAsBytes := []byte(configJson)
	configs := make([]Fofa_dic, 0)
	var jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	err := jsons.Unmarshal(jsonAsBytes, &configs)
	if err != nil {
		panic(err)
	}
	return configs
}

func get_req(url string, servers *string, headers *string, bodys *string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	*servers = r.Header.Get("Server")
	defer func() { _ = r.Body.Close() }()
	a, _ := ioutil.ReadAll(r.Body)
	*bodys = string(a)
	dataType, _ := json.Marshal(r.Header)
	*headers = string(dataType)
}

func get_server(value string, servers string) (result bool) {
	result = false
	if servers == value {
		result = true
	}
	return result
}
func get_title(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile("<title>(.*?)</title>")
	matchArr := compileRegex.FindStringSubmatch(bodys)
	if matchArr != nil {
		if matchArr[1] == value {
			result = true
		}
	}
	return result
}

func get_body(value string, bodys string) (result bool) {
	result = strings.Contains(bodys, value)
	return result
}

func get_banner(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile("(?im)<\\s*banner.*>(.*?)<\\s*/\\s*banner>")
	matchArr := compileRegex.FindAllStringSubmatch(bodys, -1)
	var i int
	for i = 0; i < len(matchArr); i++ {
		if value == matchArr[i][1] {
			result = true
		}
	}
	return result
}

func get_port(value string, urls string) (result bool) {

	u, _ := url.Parse(urls)
	result = false
	host := u.Host
	address := net.ParseIP(host)
	if address == nil {
		result = false
	} else {
		ho := strings.Split(host, ":")
		if len(ho) > 1 {
			port := ho[1]
			if port == value {
				result = true
			}
		}
	}
	return result
}

func get_head(value string, headers string) (result bool) {
	what_headers, _ := json.Marshal(headers)
	child_string := string(what_headers)
	what_headers = []byte(strings.Replace(child_string, "\"", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "\\", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "[", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "]", "", -1))
	result = strings.Contains(string(what_headers), value)
	return result
}

func check_model(match string, value string, url string, bodys string, heades string, servers string) (result bool) {
	switch {
	case match == "body_contains":
		result = get_body(value, bodys)
	case match == "protocol_contains":
		result = false
	case match == "title_contains":
		result = get_title(value, bodys)
	case match == "banner_contains":
		result = get_banner(value, bodys)
	case match == "header_contains":
		result = get_head(value, heades)
	case match == "port_contains":
		result = get_port(value, url)
	case match == "server":
		result = get_server(value, servers)
	case match == "title":
		result = get_title(value, bodys)
	case match == "cert_contains":
		result = false
	case match == "server_contains":
		result = get_server(value, servers)
	case match == "protocol":
		result = false
	default:
		result = false
	}
	return result
}
func run(url string) {
	var bodys string
	var headers string
	var servers string
	get_req(url, &servers, &headers, &bodys)
	jsonConfigList := getConfigs()
	unmarshelledConfigs := deserializeJson(jsonConfigList)
	for _, configObj := range unmarshelledConfigs {
		//fmt.Printf("Product: %s ", configObj.Rules)
		//fmt.Println(len(configObj.Rules))
		fmt.Println("正在扫描: " + configObj.Product)
		i := 0
		for i = 0; i < len(configObj.Rules); i++ {
			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, bodys, headers, servers))
				}
				if !(ContainsInSlice(arr, false)) {
					if ContainsInSlice(arr, true) {
						resulte_write(url, configObj.Product)
					}
					//fmt.Println(configObj.Rules[i][j].Match)
					//fmt.Println(configObj.Rules[i][j].Content)
				}

			} else {
				if check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, bodys, headers, servers) {
					resulte_write(url, configObj.Product)
				}
				//fmt.Println(configObj.Rules[i][0].Match)
				//fmt.Println(configObj.Rules[i][0].Content)
			}
		}
	}
}
func ContainsInSlice(items []bool, item bool) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
func resulte_write(url string, value string) {
	month := time.Now().Month()
	months := strconv.Itoa(int(month))
	day := time.Now().Day()
	days := strconv.Itoa(int(day))
	var filepath string
	filepath = months + "-" + days + ".txt"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(url + "cms :" + value + "\n")
	write.Flush()
}
func main() {
	flag.Parse()
	if *URL != "" {
		run(*URL)
	}
	if *Urllist != "" {
		file, err := os.Open(*Urllist)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		for _, line := range lines {
			go run(line)
		}
	}
	if *URL == "" && *Urllist == "" {
		fmt.Println("请使用--url= 或者 --file= 来指定目标")
		os.Exit(0)
	}
	for runtime.NumGoroutine() > 1 {
	}
	fmt.Println("done")
}
