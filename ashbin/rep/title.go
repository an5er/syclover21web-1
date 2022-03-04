package rep

import "regexp"

func Get_title(value string, bodys string) (result bool) {
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
