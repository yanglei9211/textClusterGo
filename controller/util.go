package controller

import (
	"cluster/simhash"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type BasicController struct {
	beego.Controller
}

func (b *BasicController) writeResponse(r map[string]interface{}) {
	response, err := json.Marshal(
		map[string]interface{}{
			"status": 1,
			"data":   r,
		})
	if err != nil {
		panic(err)
	}
	b.Ctx.WriteString(string(response))
}

func (b *BasicController) writeError(errInfo string) {
	response, err := json.Marshal(
		map[string]interface{}{
			"status": 0,
			"data":   errInfo,
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

var md52text map[uint64]string

func hashfunc(x string) uint64 {
	h := md5.New()
	h.Write([]byte(x))
	r := h.Sum(nil)
	var res uint64
	rs := fmt.Sprintf("%x", r[len(r)-8:])
	fmt.Sscanf(rs, "%x", &res)
	return res
}

func CacheGetter(keys Keys) (Result, error) {
	res := Result{}
	for _, k := range keys {
		sim := simhash.Simhash{}
		sim.Init(md52text[k])
		res[k] = sim.Value()
	}
	return res, nil
}

func filterTextId(ids []int, exId int) []int {
	res := make([]int, 0, len(ids))
	for _, id := range ids {
		if exId != id {
			res = append(res, id)
		}
	}
	return res
}

func getCurTime() float64 {
	t := time.Now().Unix()
	rt := float64(t)
	return rt
}

func hasTextId(db *mgo.Database, collectionName string, TextId int) (*Doc, error) {
	var doc Doc
	if notFound := db.C(collectionName).Find(bson.M{"text_id": TextId}).One(&doc); notFound != nil {
		err := fmt.Errorf(fmt.Sprintf("text_id: %d not found", TextId))
		return nil, err
	} else {
		return &doc, nil
	}
}

func emptyTextId(db *mgo.Database, collectionName string, TextId int) error {
	var doc Doc
	if notFound := db.C(collectionName).Find(bson.M{"text_id": TextId}).One(&doc); notFound == nil {
		err := fmt.Errorf(fmt.Sprintf("exists text_id: %d", TextId))
		return err
	}
	return nil
}
