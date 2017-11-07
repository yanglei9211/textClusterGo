package controller

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ClusterGetter struct {
	BasicController
}

//func (c *ClusterGetter) writeError(e error) {
//	c.writeReponse(map[string]interface{}{"error": e, "status": 0})
//}

func (c *ClusterGetter) Post() {
	input := c.Ctx.Input
	action := input.Query("action")
	var err error
	if action == "ques" {
		param := genClusterParam(c.Controller)
		res, err := Manager.QuesClusterByText(param.Tid, param.Text)
		if err != nil {
			c.writeError(err.Error())
		} else {
			c.writeReponse(map[string]interface{}{"data": res})
		}

	} else if action == "cluster" {
		param := genClusterParam(c.Controller)
		err = c.Cluster(param.Tid, param.Text)
		if err != nil {
			c.writeError(err.Error())
		} else {
			c.writeReponse(map[string]interface{}{"data": ""})
		}
	} else if action == "delete" {

	} else if action == "break" {

	} else if action == "union" {

	} else if action == "replace" {

	} else {
		c.writeError(fmt.Sprintf("no action %s", action))
	}

	//fmt.Println("post")
	//fmt.Println(action)
	//c.writeReponse(map[string]interface{}{
	//	"item": "aabbcc",
	//})
}

func (c *ClusterGetter) Cluster(TextId int, Text string) error {
	db := Manager.GetDb()
	defer db.Session.Close()
	itemSim := Manager.NewSim(Text)
	dups, err := Manager.QuesCluster(TextId, itemSim)
	if err != nil {
		return err
	}
	var selector, data bson.M
	//selector = bson.M{"text_id": TextId}
	//data = bson.M{"$set": bson.M{"text": Text}}
	//db.C("raw").Upsert(selector, data)
	if len(dups) == 0 {
		err = Manager.Add(TextId, itemSim)
		if err != nil {
			return err
		} else {
			selector = bson.M{"text_id": TextId}
			cur := time.Now().Unix()
			data = bson.M{"$set": bson.M{"rep": true, "rep_text_id": TextId, "text": Text,
				"ctime": cur, "mtime": cur, "deleted": false, "simhash": itemSim.ValueHex()}}
			db.C("data").Upsert(selector, data)
		}
	} else {
		target := dups[0]
		cur := time.Now().Unix()
		selector = bson.M{"text_id": TextId}
		data = bson.M{"$set": bson.M{"rep": false, "rep_text_id": target,
			"text": Text, "ctime": cur, "mtime": cur, "deleted": false, "simhash": itemSim.ValueHex()}}
		db.C("data").Upsert(selector, data)
	}
	return nil
}

//func (c *ClusterGetter) Post() {
//	param := genClusterParam(c.Controller)
//	res := Manager.QuesCluster(param.Text)
//	c.writeReponse(map[string]interface{}{
//		"data": res,
//	})
//}

func (c *ClusterGetter) clusterQues(edu int, quesId, quesData string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	db := Manager.GetDb()
	defer db.Session.Close()
	//var node struct{}

	//if notFound := db.C("text_info").Find(bson.M{"t_id": quesId}).Select(
	//	bson.M{"_id": 1}).One(node); notFound {
	//	fmt.Println("1")
	//} else {
	//	panic(fmt.Sprintf("t_id: %s 已经存在", quesId))
	//}
	return nil
}
