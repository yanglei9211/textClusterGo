package controller

import (
	"cluster/simhash"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"strings"
	"fmt"
	"time"
)

type Controller struct {
	Index 		*simhash.SimhashIndex
	Session 	*mgo.Session
	DBName		string
}

var Manager Controller

func (c *Controller) GetDb() *mgo.Database {
	return c.Session.Copy().DB(c.DBName)
}

func Init() {
	InitManager()
	for _, handler := range allUrls{
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
	Manager = Controller{Index:simhashIndex, Session:session, DBName:dbName}
}

func InitIndex(db *mgo.Database, collectionName string) *simhash.SimhashIndex {
	beego.Info("begin to build index")
	stTime := time.Now().UnixNano() / 1e6
	var simDoc struct{
		Id 			string 	`bson:"t_id"`
		Simhash 	string 	`bson:"simhash"`
	}

	q := db.C(collectionName).Find(bson.M{"status": "checked", "subject": "math"}).Select(
		bson.M{"simhash": 1, "t_id": 1})
	totCnt, _ := q.Count()
	beego.Info(fmt.Sprintf("index get %d docs", totCnt))

	simNodes := make([]simhash.IndexNode, totCnt)
	iter := q.Iter()
	sim := simhash.Simhash{}
	var d uint64
	for iter.Next(&simDoc) {
		fmt.Sscanf(simDoc.Simhash, "%d", &d)
		h := fmt.Sprintf("%x", d)
		sim.InitByHex(h)
		simNodes = append(simNodes, simhash.IndexNode{Sim:sim, ObjId:simDoc.Id})
	}
	resIndex := simhash.SimhashIndex{}
	resIndex.Init(simNodes)
	edTime := time.Now().UnixNano() / 1e6
	beego.Info(fmt.Sprintf("init completed cost %d ms", edTime-stTime))
	return &resIndex
}