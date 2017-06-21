package controller

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"strconv"
)

type BasicController struct {
	beego.Controller
}

func (b *BasicController) writeReponse(r map[string]interface{}) {
	response, err := json.Marshal(
		map[string]interface{}{
			"status": 1,
			"data": r,
		})
	if err != nil {
		panic(err)
	}
	b.Ctx.WriteString(string(response))
}

func bitsCount(num uint64) int {
	num = num - ((num >> 1) & 0x5555555555555555)
	num = (num & 0x3333333333333333) + ((num >> 2) & 0x3333333333333333)
	return int((((num + (num >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56)
}

func mustInt(res *int, arg string) {
	var err error
	*res, err = strconv.Atoi(arg)
	if err != nil {
		panic(err)
	}
}

func mayInt(res *int, arg string, defaultNum int) {
	if arg != "" {
		mustInt(res, arg)
	} else {
		*res = defaultNum
	}
}

func mustUnmarshal(res interface{}, arg string) {
	err := json.Unmarshal([]byte(arg), res)
	if err != nil {
		panic(err)
	}
}

func mayUnmarshal(res interface{}, arg string) {
	if arg != "" {
		mustUnmarshal(res, arg)
	}
}

func mayObjectId(res **bson.ObjectId, arg string) {
	if arg != "" && arg != "None" && arg != "null" {
		id := bson.ObjectIdHex(arg)
		*res = &id
	}
}
