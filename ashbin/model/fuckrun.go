package model

import "watch01/core"

var Bodys string
var Headers string
var Servers string
var signal chan struct{}

func Fuckrun() {
	ch := make(chan *Re_data)
	signal := make(chan struct{})
	lines := Get_urllist()
	for _, line := range lines {
		core.Wg.Add(1)
		go Get_req(line, ch)
		<-signal
		go core.Run(line, Bodys, Headers, Servers)
	}
}
