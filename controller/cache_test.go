package controller

import (
	"gopkg.in/mgo.v2"
	"testing"
	//"gopkg.in/mgo.v2/bson"
	"cluster/simhash"
	"crypto/md5"
	"fmt"
)

var ss *mgo.Database

//func getter(keys Keys) (Result, error) {
//	db := ss
//	var dt struct{
//		TextId	int	`bson:"text_id"`
//		Value   string 	`bson:"simhash"`
//	}
//	res := Result{}
//	for _, k := range keys {
//		db.C("datas").Find(bson.M{"text_id": k}).One(&dt)
//		res[dt.TextId] = dt.Value
//	}
//	return res, nil
//}
//
//func TestCache_Get(t *testing.T) {
//	ses, err := mgo.Dial("127.0.0.1:27017")
//	ss = ses.Copy().DB("cmath")
//	defer ses.Close()
//	if err != nil {
//		panic(err)
//	}
//	cc := NewCache()
//	ks := []int{0, 1, 2, 3, 4, 5, 6, 7}
//	res, err := cc.Get(ks, getter)
//	for _, k := range ks {
//		s := res[k].(string)
//		fmt.Println(k, s)
//	}
//	fmt.Println("----------")
//	ks2 := []int{0,2,4,5,8,10}
//	rs1, _ := cc.Get(ks2, getter)
//	rs2, _ := getter(ks2)
//	for _, k := range ks2 {
//		fmt.Println(k, rs1[k], rs2[k])
//	}
//}
//
//func BenchmarkCache_Get(b *testing.B) {
//	ses, err := mgo.Dial("127.0.0.1:27017")
//	ss = ses.Copy().DB("cmath")
//	defer ses.Close()
//	if err != nil {
//		panic(err)
//	}
//	cc := NewCache()
//	ks := []int{0, 1, 2, 3, 4, 12, 123, 123,4321, 13, 4123, 324, 421, 4, 4321, 4123,1234, 341}
//	for i := 0; i < b.N; i++ {
//		cc.Get(ks, getter)
//	}
//}
//
//func BenchmarkCache_Get2(b *testing.B) {
//	ses, err := mgo.Dial("127.0.0.1:27017")
//	ss = ses.Copy().DB("cmath")
//	defer ses.Close()
//	if err != nil {
//		panic(err)
//	}
//	//cc := NewCache()
//	ks := []int{0, 1, 2, 3, 4, 12, 123, 123,4321, 13, 4123, 324, 421, 4, 4321, 4123,1234, 341}
//	for i := 0; i < b.N; i++ {
//		getter(ks)
//	}
//}

// var md52text map[uint64]string

func hashfunct(x string) uint64 {
	h := md5.New()
	h.Write([]byte(x))
	r := h.Sum(nil)
	var res uint64
	rs := fmt.Sprintf("%x", r[len(r)-8:])
	fmt.Sscanf(rs, "%x", &res)
	return res
}

func getter2(keys Keys) (Result, error) {
	res := Result{}
	for _, k := range keys {
		sim := simhash.Simhash{}
		sim.Init(md52text[k])
		res[k] = sim.Value()
	}
	return res, nil
}

func TestCache_Get2(t *testing.T) {
	x := 123
	s := fmt.Sprintf("%x", x)
	fmt.Println(s)
	ss := []string{
		"关于函数y-dfrac1x的图象下列说法中错误的是leftqquadright",
		"二次函数yax^2+bx+c的图象如图反比例函数ydfracax与正比例函数ybx在同一坐标系内的大致图象是leftqquadright",
		"如图在triangleABC和triangleDEC中已知ABDE还需要添加两个条件才能使triangleABCcongtriangleDEC不能添加的一组是leftqquadright",
		"在-3-201mathrmpi这五个数中最小的实数是leftqquadright",
		"下面四个几何体中左视图是四边形的几何体共有leftqquadright",
		"sqrt4的算术平方根是leftqquadright",
		"小明在学习了正方形之后给同桌小文出了道题从下列四个条件ABBCangleABC90^circACBDACperpBD中选两个作为补充条件使平行四边形ABCD成为正方形如图现有下列四种选法你认为其中错误的是leftqquadright",
		"正方形网格中angleAOB如下图放置则sin angleAOBleftqquadright",
		"在triangleABC中angleC90^circD是AC上的一点DEperpAB于点E若AC8BC6AD5则DE的长为leftqquadright",
		"课题研究小组对附着在物体表面的三个微生物课题小组成员把它们分别标号为123的生长情况进行观察记录这三个微生物第一天各自一分为二产生新的微生物分别被标号为456789接下去每天都按照这样的规律变化即每个微生物分为二形成新的微生物课题组成员用如图所示的图形进行形象的记录那么标号为100的微生物会出现在leftqquadright",
		"如图EFparallelBCAC平分angleBAFangleB80^circ则angleC的度数为leftqquadright",
	}
	md52text = make(map[uint64]string)
	keys := make([]uint64, 0, len(ss))
	for _, s := range ss {
		m := hashfunc(s)
		md52text[m] = s
		keys = append(keys, m)
	}
	cc := NewCache()
	res, _ := cc.Get(keys, getter2)
	fmt.Println("-----------------------")
	fmt.Println(res)
}

func BenchmarkCache_Get3(b *testing.B) {
	ss := []string{
		"关于函数y-dfrac1x的图象下列说法中错误的是leftqquadright",
		"二次函数yax^2+bx+c的图象如图反比例函数ydfracax与正比例函数ybx在同一坐标系内的大致图象是leftqquadright",
		"如图在triangleABC和triangleDEC中已知ABDE还需要添加两个条件才能使triangleABCcongtriangleDEC不能添加的一组是leftqquadright",
		"在-3-201mathrmpi这五个数中最小的实数是leftqquadright",
		"下面四个几何体中左视图是四边形的几何体共有leftqquadright",
		"sqrt4的算术平方根是leftqquadright",
		"小明在学习了正方形之后给同桌小文出了道题从下列四个条件ABBCangleABC90^circACBDACperpBD中选两个作为补充条件使平行四边形ABCD成为正方形如图现有下列四种选法你认为其中错误的是leftqquadright",
		"正方形网格中angleAOB如下图放置则sin angleAOBleftqquadright",
		"在triangleABC中angleC90^circD是AC上的一点DEperpAB于点E若AC8BC6AD5则DE的长为leftqquadright",
		"课题研究小组对附着在物体表面的三个微生物课题小组成员把它们分别标号为123的生长情况进行观察记录这三个微生物第一天各自一分为二产生新的微生物分别被标号为456789接下去每天都按照这样的规律变化即每个微生物分为二形成新的微生物课题组成员用如图所示的图形进行形象的记录那么标号为100的微生物会出现在leftqquadright",
		"如图EFparallelBCAC平分angleBAFangleB80^circ则angleC的度数为leftqquadright",
	}
	md52text = make(map[uint64]string)
	keys := make([]uint64, 0, len(ss))
	for _, s := range ss {
		m := hashfunc(s)
		md52text[m] = s
		keys = append(keys, m)
	}
	for i := 0; i < b.N; i++ {
		getter2(keys)
	}
}

func BenchmarkCache_Get4(b *testing.B) {
	ss := []string{
		"关于函数y-dfrac1x的图象下列说法中错误的是leftqquadright",
		"二次函数yax^2+bx+c的图象如图反比例函数ydfracax与正比例函数ybx在同一坐标系内的大致图象是leftqquadright",
		"如图在triangleABC和triangleDEC中已知ABDE还需要添加两个条件才能使triangleABCcongtriangleDEC不能添加的一组是leftqquadright",
		"在-3-201mathrmpi这五个数中最小的实数是leftqquadright",
		"下面四个几何体中左视图是四边形的几何体共有leftqquadright",
		"sqrt4的算术平方根是leftqquadright",
		"小明在学习了正方形之后给同桌小文出了道题从下列四个条件ABBCangleABC90^circACBDACperpBD中选两个作为补充条件使平行四边形ABCD成为正方形如图现有下列四种选法你认为其中错误的是leftqquadright",
		"正方形网格中angleAOB如下图放置则sin angleAOBleftqquadright",
		"在triangleABC中angleC90^circD是AC上的一点DEperpAB于点E若AC8BC6AD5则DE的长为leftqquadright",
		"课题研究小组对附着在物体表面的三个微生物课题小组成员把它们分别标号为123的生长情况进行观察记录这三个微生物第一天各自一分为二产生新的微生物分别被标号为456789接下去每天都按照这样的规律变化即每个微生物分为二形成新的微生物课题组成员用如图所示的图形进行形象的记录那么标号为100的微生物会出现在leftqquadright",
		"如图EFparallelBCAC平分angleBAFangleB80^circ则angleC的度数为leftqquadright",
	}
	md52text = make(map[uint64]string)
	keys := make([]uint64, 0, len(ss))
	for _, s := range ss {
		m := hashfunc(s)
		md52text[m] = s
		keys = append(keys, m)
	}
	cc := NewCache()
	for i := 0; i < b.N; i++ {
		cc.Get(keys, getter2)
	}
}
