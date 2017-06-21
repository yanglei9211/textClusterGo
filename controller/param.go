package controller

import (

)

type BaseParam struct {
	quesId 			string
	quesData		string
}

type ClusterParam struct {
	BaseParam
}

type AddParam struct {
	BaseParam
}

type DeleteParam struct {
	quesId			string
}

type UnionParam struct {
	repId			string
	dupId			string
}

type SeparateParam	struct {
	quesId 			string
}

type AutoUnionParam	struct {
	BaseParam
	cTime			float64
}

type ReplaceRepParam struct {
	oldRepId		string
	newRepId		string
}

type CalcParam struct {

}