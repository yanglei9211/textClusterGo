package controller

import (
	"fmt"
)

type ClusterGetter struct {
	BasicController
}

func (c *ClusterGetter) Post() {
	input := c.Ctx.Input
	action := input.Query("action")
	fmt.Println("post")
	fmt.Println(action)
	c.writeReponse(map[string]interface{}{
		"item": "aabbcc",
	})
}

func (c *ClusterGetter) clusterQues(edu int, quesId, quesData string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	db := Manager.GetDb()
	defer db.Session.Close()
	var node struct{}

	if notFound := db.C("text_info").Find(bson.M{"t_id": quesId}).Select(
		bson.M{"_id": 1}).One(node); notFound {

	} else {
		panic(fmt.Sprintf("t_id: %s 已经存在", quesId))
	}
	return nil
}
