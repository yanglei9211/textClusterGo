package controller

import (
	"cluster/simhash"
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"strings"
	"time"
)

type Controller struct {
	// 应该是bl吧
	Index   *simhash.SimhashIndex
	Session *mgo.Session
	DBName  string
	Cache   *Cache
}

type Doc struct {
	TextId    int    `bson:"text_id"`
	Rep       bool   `bson:"rep"`
	RepTextId int    `bson:"rep_text_id"`
	Deleted   bool   `bson:"deleted"`
	Text      string `bson:"text"`
	Ctime     int64  `bson:"ctime"`
	Mtime     int64  `bson:"mtime"`
}

func (c *Controller) GetDb() *mgo.Database {
	return c.Session.Copy().DB(c.DBName)
}

func (c *Controller) QuesClusterByText(TextId int, Text string) (sort.IntSlice, error) {
	s := c.NewSim(Text)
	return c.QuesCluster(TextId, s)
}

func (c *Controller) AddBk(TextId int, Text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	tpNode := NewIndexNode(c.NewSim(Text), TextId)
	c.Index.Add(*tpNode)
	return nil
}

func (c *Controller) QuesCluster(TextId int, s *simhash.Simhash) (res sort.IntSlice, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	dups := c.Index.GetNearDups(*s)
	res = make(sort.IntSlice, len(dups))
	for i, s := range dups {
		fmt.Sscanf(s, "%d", &res[i])
	}
	res.Sort()
	res = filterTextId(res, TextId)
	return res, nil
}

func (c *Controller) Add(TextId int, s *simhash.Simhash) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	tpNode := NewIndexNode(s, TextId)
	c.Index.Add(*tpNode)
	return nil
}

func (c *Controller) Del(TextId int, Text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	sim := c.NewSim(Text)
	dups := c.Index.GetNearDups(*sim)
	if len(dups) == 0 {
		panic("被删除元素不存在")
	}
	tpNode := NewIndexNode(sim, TextId)
	c.Index.Del(*tpNode)
	return nil
}

func (c *Controller) Cluster(TextId, RepId int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	db := c.GetDb()
	db.Session.Close()
	selector := bson.M{"text_id": TextId}
	data := bson.M{"rep_id": RepId}
	db.C("doc").Update(selector, data)
	return nil
}

func (c *Controller) NewSim(Text string) *simhash.Simhash {
	newSimhash := func(text string) *simhash.Simhash {
		sim := simhash.Simhash{}
		sim.Init(text)
		return &sim
	}

	md52text = make(map[uint64]string)
	m := hashfunc(Text)
	md52text[m] = Text
	keys := Keys{m}
	res, err := c.Cache.Get(keys, CacheGetter)
	if err != nil {
		fmt.Println("cache error!!!!!!")
		res := newSimhash(Text)
		return res
	} else {
		fmt.Println("use cache!!!!!!!")
		res := NewSimhashByVal(res[m].(uint64))
		return res
	}
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
	Manager.Cache = NewCache()
}

func InitIndex(db *mgo.Database, collectionName string) *simhash.SimhashIndex {
	beego.Info("begin to build index")
	stTime := time.Now().UnixNano() / 1e6
	var simDoc struct {
		TextId  string `bson:"text_id"`
		Simhash string `bson:"simhash"`
	}

	q := db.C(collectionName).Find(bson.M{"rep": true}).Select(
		bson.M{"simhash": 1, "text_id": 1})
	totCnt, _ := q.Count()
	beego.Info(fmt.Sprintf("index get %d docs", totCnt))
	simNodes := make([]simhash.IndexNode, 0, totCnt)
	iter := q.Iter()
	for iter.Next(&simDoc) {
		simNodes = append(simNodes, simhash.IndexNode{Sim: *NewSimhashByHex(simDoc.Simhash), ObjId: simDoc.TextId})
	}
	resIndex := simhash.SimhashIndex{}
	resIndex.Init(simNodes)
	edTime := time.Now().UnixNano() / 1e6
	beego.Info(fmt.Sprintf("init completed cost %d ms", edTime-stTime))
	fmt.Println("okokok")
	return &resIndex
}

func NewSimhashByHex(h string) *simhash.Simhash {
	sim := simhash.Simhash{}
	sim.InitByHex(h)
	return &sim
}

func NewSimhashByVal(v uint64) *simhash.Simhash {
	sim := simhash.Simhash{}
	sim.InitByValue(v)
	return &sim
}

func NewIndexNode(sim *simhash.Simhash, TextId int) *simhash.IndexNode {
	ObjId := fmt.Sprintf("%d", TextId)
	node := simhash.IndexNode{}
	node.Init(*sim, ObjId)
	return &node
}
