package rep

import "regexp"

func Get_banner(value string, bodys string) (result bool) {
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
