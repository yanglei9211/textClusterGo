package controller

import (
	"fmt"
)

type TestGetter struct {
	BasicController
}

func (sg *TestGetter) Get() {
	fmt.Println("get!!!!")
	res := "aaabbbccc"
	sg.writeResponse(map[string]interface{}{
		"item": res,
	})
}

func (sg *TestGetter) Post() {
	fmt.Println("post!!!!!")
	res := "cccbbbaaa"
	sg.writeResponse(map[string]interface{}{
		"item": res,
	})
}

func (sg *TestGetter) Put() {
	fmt.Println("put!!!!!")
	res := "bbbcccaaa"
	sg.writeResponse(map[string]interface{}{
		"item": res,
	})
}

func (sg *TestGetter) Delete() {
	fmt.Println("delete!!!!!")
	res := "cccaaabbb"
	sg.writeResponse(map[string]interface{}{
		"item": res,
	})
}
