package controller

import (
	"cluster/simhash"
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Controller struct {
	Index   *simhash.SimhashIndex
	Session *mgo.Session
	DBName  string
}

func (c *Controller) GetDb() *mgo.Database {
	return c.Session.Copy().DB(c.DBName)
}

func (c *Controller) QuesCluster(Text string) (res []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	res = c.Index.GetNearDups(NewSimhash(Text))
	return res, nil
}

func (c *Controller) Add(TextId, Text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	tpNode := NewIndexNode(NewSimhash(Text), TextId)
	c.Index.Add(tpNode)
	return nil
}

func (c *Controller) Del(TextId, Text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	sim := NewSimhash(Text)
	dups := c.Index.GetNearDups(sim)
	if len(dups) == 0 {
		panic("被删除元素不存在")
	}
	tpNode := NewIndexNode(sim, TextId)
	c.Index.Del(tpNode)
	return nil
}

var Manager Controller

func Init() {
	InitManager()
	for _, handler := range allUrls {
		beego.Router(handler.url, handler.controller)
	}
}

func InitManager() {
	appConfig := beego.AppConfig
	dbHost := appConfig.String("db_host")
	session, err := mgo.Dial(dbHost)
	if err != nil {
		panic(err)
	}
	collectionName := strings.TrimSpace(appConfig.String("collection_name"))
	dbName := strings.TrimSpace(appConfig.String("db_name"))
	db := session.DB(dbName)
	simhashIndex := InitIndex(db, collectionName)
	Manager = Controller{Index: simhashIndex, Session: session, DBName: dbName}
}

func InitIndex(db *mgo.Database, collectionName string) *simhash.SimhashIndex {
	beego.Info("begin to build index")
	stTime := time.Now().UnixNano() / 1e6
	var simDoc struct {
		TextId  string    `bson:"text_id"`
		Simhash string    `bson:"simhash"`
	}

	q := db.C(collectionName).Find(bson.M{}).Select(
		bson.M{"simhash": 1, "text_id": 1})
	totCnt, _ := q.Count()
	beego.Info(fmt.Sprintf("index get %d docs", totCnt))
	simNodes := make([]simhash.IndexNode, 0, totCnt)
	iter := q.Iter()
	for iter.Next(&simDoc) {
		simNodes = append(simNodes, simhash.IndexNode{Sim: NewSimhashByHew(simDoc.Simhash), ObjId: simDoc.TextId})
	}
	resIndex := simhash.SimhashIndex{}
	resIndex.Init(simNodes)
	edTime := time.Now().UnixNano() / 1e6
	beego.Info(fmt.Sprintf("init completed cost %d ms", edTime-stTime))
	fmt.Println("okokok")
	return &resIndex
}

func NewSimhash(data string) simhash.Simhash {
	sim := simhash.Simhash{}
	sim.Init(data)
	return sim
}

func NewSimhashByHew(h string) simhash.Simhash {
	sim := simhash.Simhash{}
	sim.InitByHex(h)
	return sim
}

func NewIndexNode(sim simhash.Simhash, ObjId string) simhash.IndexNode {
	node := simhash.IndexNode{}
	node.Init(sim, ObjId)
	return node
}