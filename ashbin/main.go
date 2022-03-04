package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"watch01/core"
	"watch01/model"
)

func main() {
	ch := make(chan string)
	core.Get_pa()
	if *core.URL != "" {
		model.Get_req(*core.URL, ch)
		core.Run(*core.URL, model.Bodys, model.Headers, model.Servers)
	}
	if *core.Urllist != "" {
		p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
			model.Fuckrun()
			core.Wg.Done()
		})
		defer p.Release()
	}
	if *core.URL == "" && *core.Urllist == "" {
		fmt.Println("请使用--url= 或者 --file= 来指定目标")
		os.Exit(0)
	}
	core.Wg.Wait()
	fmt.Println(*core.URL)
	fmt.Println(*core.Urllist)
	fmt.Println(&model.Bodys)
	fmt.Println(core.Wg)
	fmt.Println("done")
}
