package rep

import "strings"

func Get_body(value string, bodys string) (result bool) {
	result = strings.Contains(bodys, value)
	return result
}
