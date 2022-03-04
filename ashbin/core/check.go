package core

import "watch01/rep"

func Check_model(match string, value string, url string, bodys string, heades string, servers string) (result bool) {
	switch {
	case match == "body_contains":
		result = rep.Get_body(value, bodys)
	case match == "protocol_contains":
		result = false
	case match == "title_contains":
		result = rep.Get_title(value, bodys)
	case match == "banner_contains":
		result = rep.Get_banner(value, bodys)
	case match == "header_contains":
		result = rep.Get_head(value, heades)
	case match == "port_contains":
		result = rep.Get_port(value, url)
	case match == "server":
		result = rep.Get_server(value, servers)
	case match == "title":
		result = rep.Get_title(value, bodys)
	case match == "cert_contains":
		result = false
	case match == "server_contains":
		result = rep.Get_server(value, servers)
	case match == "protocol":
		result = false
	default:
		result = false
	}
	return result
}
