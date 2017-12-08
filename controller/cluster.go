package controller

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type ClusterGetter struct {
	BasicController
}

//func (c *ClusterGetter) writeError(e error) {
//	c.writeResponse(map[string]interface{}{"error": e, "status": 0})
//}

func (c *ClusterGetter) Post() {
	input := c.Ctx.Input
	action := input.Query("action")
	var err error
	if action == "ques" {
		param := genAutoClusterParam(c.Controller)
		res, err := Manager.QuesClusterByText(param.Tid, param.Text)
		if err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": res})
		}
	} else if action == "add" {
		param := genAddParam(c.Controller)
		if err = c.AddDoc(param.Tid, param.Text); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}
	} else if action == "cluster" {
		param := genClusterParam(c.Controller)
		if err = c.Cluster(param.Cid, param.Tid, param.Text); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}

	} else if action == "auto_cluster" {
		param := genAutoClusterParam(c.Controller)
		if tarId, err := c.AutoCluster(param.Tid, param.Text); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": tarId})
		}
	} else if action == "delete" {
		param := genDeleteParam(c.Controller)
		if err = c.DeleteDoc(param.Tid); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}

	} else if action == "break" {
		param := genBreakParam(c.Controller)
		if err = c.Break(param.Tid, param.Cid); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}
	} else if action == "union" {
		param := genUnionParam(c.Controller)
		if err = c.Union(param.Tid, param.Cid); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}

	} else if action == "replace" {
		param := genReplaceParam(c.Controller)
		if err = c.Replace(param.Tid); err != nil {
			c.writeError(err.Error())
		} else {
			c.writeResponse(map[string]interface{}{"data": ""})
		}
	} else {
		c.writeError(fmt.Sprintf("no action %s", action))
	}

	//fmt.Println("post")
	//fmt.Println(action)
	//c.writeResponse(map[string]interface{}{
	//	"item": "aabbcc",
	//})
}

func (c *ClusterGetter) AddDoc(TextId int, Text string) error {
	sim := Manager.NewSim(Text)
	if err := Manager.Add(TextId, sim); err != nil {
		return err
	}
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName

	if idExists := emptyTextId(db, collection_name, TextId); idExists != nil {
		return idExists
	}
	cur := getCurTime()
	selector := bson.M{"text_id": TextId}
	data := bson.M{"$set": bson.M{"rep": true, "rep_text_id": TextId, "text": Text,
		"ctime": cur, "mtime": cur, "deleted": false, "simhash": sim.ValueHex()}}
	db.C(collection_name).Upsert(selector, data)
	return nil
}

func (c *ClusterGetter) Cluster(ClusId, TextId int, Text string) error {
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	if idExists := emptyTextId(db, collection_name, TextId); idExists != nil {
		return idExists
	}
	if clusDoc, found := hasTextId(db, collection_name, ClusId); found != nil || !clusDoc.Rep {
		return fmt.Errorf(fmt.Sprintf("目标text_id: %d状态不符合或不存在", ClusId))
	}
	sim := Manager.NewSim(Text)
	cur := getCurTime()
	selector := bson.M{"text_id": TextId}
	data := bson.M{"$set": bson.M{"rep": false, "rep_text_id": ClusId, "text": Text,
		"ctime": cur, "mtime": cur, "deleted": false, "simhash": sim.ValueHex()}}
	db.C(collection_name).Upsert(selector, data)
	return nil
}

func (c *ClusterGetter) AutoCluster(TextId int, Text string) (int, error) {
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	if idExists := emptyTextId(db, collection_name, TextId); idExists != nil {
		return -1, idExists
	}
	itemSim := Manager.NewSim(Text)
	dups, err := Manager.QuesCluster(TextId, itemSim)
	if err != nil {
		return -1, err
	}
	var selector, data bson.M
	//selector = bson.M{"text_id": TextId}
	//data = bson.M{"$set": bson.M{"text": Text}}
	//db.C("raw").Upsert(selector, data)
	var tarId int
	if len(dups) == 0 {
		tarId = TextId
		err = Manager.Add(TextId, itemSim)
		if err != nil {
			return -1, err
		} else {
			selector = bson.M{"text_id": TextId}
			//cur := time.Now().Unix()
			cur := getCurTime()
			data = bson.M{"$set": bson.M{"rep": true, "rep_text_id": TextId, "text": Text,
				"ctime": cur, "mtime": cur, "deleted": false, "simhash": itemSim.ValueHex()}}
			db.C(collection_name).Upsert(selector, data)
		}
	} else {
		tarId = dups[0]
		//cur := time.Now().Unix()
		cur := getCurTime()
		selector = bson.M{"text_id": TextId}
		data = bson.M{"$set": bson.M{"rep": false, "rep_text_id": tarId,
			"text": Text, "ctime": cur, "mtime": cur, "deleted": false, "simhash": itemSim.ValueHex()}}
		db.C(collection_name).Upsert(selector, data)
	}
	return tarId, nil
}

func (c *ClusterGetter) DeleteDoc(TextId int) error {
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	var doc Doc
	if notFound := db.C(collection_name).Find(bson.M{"text_id": TextId, "deleted": false}).One(&doc); notFound != nil {
		return notFound
		//panic(fmt.Sprintf("text_id: %s 不存在", TextId))
	} else {
		if doc.Rep {
			// 如果是代表,连同dups全部删除
			if err := Manager.DelWithSim(TextId, doc.Sim); err == nil {
				selector := bson.M{"rep_text_id": TextId}
				data := bson.M{"$set": bson.M{"deleted": true, "mtime": getCurTime()}}
				db.C(collection_name).UpdateAll(selector, data)
			} else {
				return err
			}
		} else {
			// 否则,只删除目标文档
			selector := bson.M{"text_id": TextId}
			data := bson.M{"$set": bson.M{"deleted": true, "mtime": getCurTime()}}
			db.C(collection_name).Update(selector, data)
		}
	}
	return nil
}

func (c *ClusterGetter) Break(TextId, ClusId int) error {
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	var doc Doc
	if notFound := db.C(collection_name).Find(bson.M{"text_id": TextId, "deleted": false,
		"rep": false, "rep_text_id": ClusId}).One(&doc); notFound != nil {
		return notFound
		//panic(fmt.Sprintf("text_id: %s 状态不符合", TextId))
	} else {
		if err := Manager.AddWithSim(TextId, doc.Sim); err == nil {
			cur := getCurTime()
			selector := bson.M{"text_id": TextId}
			data := bson.M{"$set": bson.M{"rep": true, "rep_text_id": TextId, "mtime": cur}}
			db.C(collection_name).Update(selector, data)
		} else {
			return err
		}
	}
	return nil
}

func (c *ClusterGetter) Union(TextId, ClusId int) error {
	/*cluster textid to clusid*/
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	var doc, clusDoc Doc
	var notFound error
	if notFound = db.C(collection_name).Find(bson.M{"text_id": TextId, "rep": true, "deleted": false}).One(&doc); notFound != nil {
		return fmt.Errorf(fmt.Sprintf("被合并text_id: %d 状态不符合或不存在", TextId))
	}
	if notFound = db.C(collection_name).Find(bson.M{"text_id": ClusId, "rep": true, "deleted": false}).One(&clusDoc); notFound != nil {
		return fmt.Errorf(fmt.Sprintf("目标text_id: %d 状态不符合或不存在", ClusId))
	}
	if err := Manager.DelWithSim(TextId, doc.Sim); err == nil {
		cur := getCurTime()
		selector := bson.M{"text_id": TextId}
		data := bson.M{"rep": false, "rep_text_id": ClusId, "mtime": cur}
		db.C(collection_name).Update(selector, data)
	}
	return nil
}

func (c *ClusterGetter) Replace(TextId int) error {
	db := Manager.GetDb()
	defer db.Session.Close()
	collection_name := Manager.CollectionName
	var doc, clusDoc Doc
	var notFound error
	if notFound = db.C(collection_name).Find(bson.M{"text_id": TextId, "rep": false}).One(&doc); notFound != nil {
		err := fmt.Errorf("text_id: %d状态不符合", TextId)
		return err
	}
	if notFound = db.C(collection_name).Find(bson.M{"text_id": doc.RepTextId, "rep": true}).One(&clusDoc); notFound != nil {
		err := fmt.Errorf("text_id: %d状态不符合", doc.RepTextId)
		return err
	}
	if errDel, errAdd := Manager.DelWithSim(clusDoc.TextId, clusDoc.Sim), Manager.AddWithSim(doc.TextId, doc.Sim); errDel == nil && errAdd == nil {
		cur := getCurTime()
		selector := bson.M{"rep_text_id": clusDoc.TextId}
		data := bson.M{"$set": bson.M{"rep_text_id": TextId, "mtime": cur}}
		db.C(collection_name).UpdateAll(selector, data)

		selector = bson.M{"text_id": clusDoc.TextId}
		data = bson.M{"$set": bson.M{"rep": false}}
		db.C(collection_name).Update(selector, data)
		selector = bson.M{"text_id": TextId}
		data = bson.M{"$set": bson.M{"rep": true}}
		db.C(collection_name).Update(selector, data)
		//selector = bson.M{"text_id": clusDoc.TextId}
		//data = bson.M{"rep": false, "rep_text_id": doc.TextId, "mtime": cur}
		//db.C(collection_name).Update(selector, data)
		//selector = bson.M{"text_id": doc.Text}
		//data = bson.M{"rep": true, "rep_Text_id": doc.TextId, "mtime": cur}
		//db.C(collection_name).Update(selector, data)

	} else {
		if errDel != nil {
			return errDel
		} else {
			return errAdd
		}
	}

	return nil
}

//func (c *ClusterGetter) Post() {
//	param := genClusterParam(c.Controller)
//	res := Manager.QuesCluster(param.Text)
//	c.writeResponse(map[string]interface{}{
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
