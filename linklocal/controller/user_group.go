package controller

import (
	"gopkg.in/macaron.v1"
	"linklocal/auth"
	"linklocal/utils"
	"fmt"
)

func GetGroupBindingUsers(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	users,err := auth.GetGroupBindingUsers(groupId)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["users"] = make([]map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["users"] = users
	ctx.JSON(200,info)
	return
}

func GetUserBindingGroups(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	groups,err := auth.GetUserBindingGroups(userId)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["groups"] = make([]map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["groups"] = groups
	ctx.JSON(200,info)
	return
}

func BindUserToGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetUserGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.BindUserToGroup(groupId,args)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	ctx.JSON(200,info)
	return
}

func UnbindUserToGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetUserGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	_,ok1 := args["user_id"]
	if !ok1{
		err = fmt.Errorf("UserId or groupId not set")
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.UnbindUserToGroup(groupId,args["user_id"].(string))
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	ctx.JSON(200,info)
	return
}

func SetUserTypeForGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetUserGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	_,ok1 := args["type"]
	_,ok2 := args["user_id"]
	if !(ok1 && ok2){
		err = fmt.Errorf("UserId or Type not set")
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.SetUserTypeForGroup(groupId,args["user_id"].(string),args["type"].(string))
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	ctx.JSON(200,info)
	return
}

func SetUserPriorityForGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetUserGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	_,ok1 := args["priority"]
	_,ok2 := args["user_id"]
	if !(ok1 && ok2){
		err = fmt.Errorf("UserId or priority not set")
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.SetUserPriorityForGroup(groupId,args["user_id"].(string),args["priority"].(string))
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	ctx.JSON(200,info)
	return
}

//For Get POST Info Only
func GetUserGroupPostInfo(ctx *macaron.Context)(args map[string]interface{}, err error){
	body, err := utils.GetPostBody(ctx)
	tokenStr := ctx.Req.Header.Get("token")
	user, err := auth.GetUserByTokenString(tokenStr)
	if err != nil{
		return
	}
	args = make(map[string]interface{},0)
	for k,v := range body {
		switch k {
		case "user_id":
			args["user_id"] = v.(string)
		case "group_id":
			args["group_id"] = v.(string)
		case "type":
			args["type"] = v.(string)
		case "priority":
			args["priority"] = v.(string)
		case "end_time":
			args["end_time"] = utils.InterfaceToInt(v)
		}
	}
	args["current_user"] = user.Name
	return
}
