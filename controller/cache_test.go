package controller

import (
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

var ss *mgo.Database

func getter(keys Keys) (Result, error) {
	db := ss
	var dt struct{
		TextId	int	`bson:"text_id"`
		Value   string 	`bson:"simhash"`
	}
	res := Result{}
	for _, k := range keys {
		db.C("datas").Find(bson.M{"text_id": k}).One(&dt)
		res[dt.TextId] = dt.Value
	}
	return res, nil
}

func TestCache_Get(t *testing.T) {

	ses, err := mgo.Dial("127.0.0.1:27017")
	ss = ses.Copy().DB("cmath")
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	cc := NewCache()
	ks := []int{0, 1, 2, 3, 4, 5, 6, 7}
	res, err := cc.Get(ks, getter)
	for _, k := range ks {
		s := res[k].(string)
		fmt.Println(k, s)
	}
	fmt.Println("----------")
	ks2 := []int{0,2,4,5,8,10}
	rs1, _ := cc.Get(ks2, getter)
	rs2, _ := getter(ks2)
	for _, k := range ks2 {
		fmt.Println(k, rs1[k], rs2[k])
	}
}

func BenchmarkCache_Get(b *testing.B) {
	ses, err := mgo.Dial("127.0.0.1:27017")
	ss = ses.Copy().DB("cmath")
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	cc := NewCache()
	ks := []int{0, 1, 2, 3, 4, 12, 123, 123,4321, 13, 4123, 324, 421, 4, 4321, 4123,1234, 341}
	for i := 0; i < b.N; i++ {
		cc.Get(ks, getter)
	}
}

func BenchmarkCache_Get2(b *testing.B) {
	ses, err := mgo.Dial("127.0.0.1:27017")
	ss = ses.Copy().DB("cmath")
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	//cc := NewCache()
	ks := []int{0, 1, 2, 3, 4, 12, 123, 123,4321, 13, 4123, 324, 421, 4, 4321, 4123,1234, 341}
	for i := 0; i < b.N; i++ {
		getter(ks)
	}
}