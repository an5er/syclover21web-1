package rep

func Get_server(value string, servers string) (result bool) {
	result = false
	if servers == value {
		result = true
	}
	return result
}
