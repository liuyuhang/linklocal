package utils

import (
	"encoding/json"
	"gopkg.in/macaron.v1"
)

//Get Post Body args
func GetPostBody(ctx *macaron.Context)(m map[string]interface{}, err error){
	bodyInByte,err := ctx.Req.Body().Bytes()
	if err != nil{
		return
	}
	// var f interface{}
	err = json.Unmarshal(bodyInByte, &m)
	return
}
