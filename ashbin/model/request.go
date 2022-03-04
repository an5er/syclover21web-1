package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Re_data struct {
	url    string
	server string
	header string
	bodys  string
}

func Get_req(url string, ch chan *Re_data) {
	var re_data Re_data
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	re_data.server = r.Header.Get("Server")
	defer func() { _ = r.Body.Close() }()
	a, _ := ioutil.ReadAll(r.Body)
	re_data.bodys = string(a)
	dataType, _ := json.Marshal(r.Header)
	re_data.header = string(dataType)
	re_data.url = url
	var pointer = &re_data
	ch <- pointer
}
