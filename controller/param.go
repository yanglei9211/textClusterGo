package controller

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type Param struct {
	Action string
	Tid    int
	Text   string
	Cid    int
}

func genParam(c beego.Controller) *Param {
	defer func() {
		if err := recover(); err != nil {
			c.Ctx.Output.SetStatus(400)
		}
	}()
	input := c.Ctx.Input
	action := input.Query("action")
	param := new(Param)
	if action == "cluster" {
		param = genClusterParam(c)
	} else if action == "break" {
		param = genBreakParam(c)
	} else if action == "ques" {
		param = genQuesParam(c)
	} else if action == "union" {
		param = genUnionParam(c)
	} else {
		panic(errors.New(fmt.Sprintf("no such action: %s", action)))
	}
	return param
}

func genClusterParam(c beego.Controller) *Param {
	input := c.Ctx.Input
	cparam := new(Param)
	mustInt(&cparam.Tid, input.Query("text_id"))
	cparam.Text = input.Query("text")
	return cparam
}

func genBreakParam(c beego.Controller) *Param {
	input := c.Ctx.Input
	bparam := new(Param)
	mustInt(&bparam.Tid, input.Query("text_id"))
	mustInt(&bparam.Cid, input.Query("clus_id"))
	return bparam
}

func genQuesParam(c beego.Controller) *Param {
	input := c.Ctx.Input
	qparam := new(Param)
	mustInt(&qparam.Tid, input.Query("text_id"))
	qparam.Text = input.Query("text")
	return qparam
}

func genUnionParam(c beego.Controller) *Param {
	input := c.Ctx.Input
	uparam := new(Param)
	mustInt(&uparam.Tid, input.Query("text_id"))
	mustInt(&uparam.Cid, input.Query("clus_id"))
	return uparam
}

func genReplaceParam(c beego.Controller) *Param {
	input := c.Ctx.Input
	uparam := new(Param)
	mustInt(&uparam.Tid, input.Query("text_id"))
	return uparam
}
