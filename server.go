package main

import (
	"github.com/astaxie/beego"
	"runtime"
	"cluster/controller"
)

func main() {

	controller.Init()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	beego.Run()
	/*
	ss := []string{"放大是否打算","发生的旅客合法的顺口溜合法的索科洛夫哈萨克了会发生", "11223344", "aabbcc",
		"AABBCC", "AaBbcC"}
	res := []simhash.Simhash{}
	for _, s := range ss {
		sim := simhash.Simhash{}
		sim.Init(s)
		res = append(res, sim)
	}
	for _, r := range res {
		fmt.Println(r.Value())
	}
	test := []simhash.IndexNode{}
	s := simhash.SimhashIndex{}
	for idx := range res {
		tpNode := simhash.IndexNode{res[idx], fmt.Sprintf("%d", idx)}
		test = append(test, tpNode)
	}
	s.Init(test)
	ts := simhash.Simhash{}
	ts.Init("aabbcc")
	fmt.Println(ts.Value())
	ans := s.GetNearDups(ts)
	fmt.Println(ans)
	*/
}
