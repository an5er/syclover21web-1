package rep

import (
	"net"
	"net/url"
	"strings"
)

func Get_port(value string, urls string) (result bool) {

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
