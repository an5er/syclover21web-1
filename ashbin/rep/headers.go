package rep

import (
	"encoding/json"
	"strings"
)

func Get_head(value string, headers string) (result bool) {
	what_headers, _ := json.Marshal(headers)
	child_string := string(what_headers)
	what_headers = []byte(strings.Replace(child_string, "\"", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "\\", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "[", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "]", "", -1))
	result = strings.Contains(string(what_headers), value)
	return result
}
