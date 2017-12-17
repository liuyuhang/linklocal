package controller

import (
	"gopkg.in/macaron.v1"
	"linklocal/auth"
	"linklocal/utils"
	"linklocal/model"
)

func CreateGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	args,err := GetGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	group, err := auth.CreateGroup(args)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["group"] = group
	ctx.JSON(200,info)
	return
}

func DeleteGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	err := auth.DeleteGroup(groupId)
	if err != nil{
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

func UpdateGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	group,err := auth.UpdateGroup(groupId,args)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["group"] = group
	ctx.JSON(200,info)
	return

}

func ListGroups(ctx *macaron.Context){
	info := make(map[string]interface{})
	groups,err := auth.ListGroups()
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["groups"] = make([]model.Group,0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["groups"] = groups
	ctx.JSON(200,info)
	return
}

func GetGroup(ctx *macaron.Context){
	info := make(map[string]interface{})
	id := ctx.Params(":id")
	group,err := auth.GetGroup(id)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{})
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["group"] = group
	ctx.JSON(200,info)
	return
}

func ChangeGroupEndTime(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	args,err := GetGroupPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.ChangeGroupEndTime(groupId,args)
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

func GetGroupParent(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	group,err := auth.GetGroupParent(groupId)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["group"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["group"] = group
	ctx.JSON(200,info)
	return
}

func GetGroupChildren(ctx *macaron.Context){
	info := make(map[string]interface{})
	groupId := ctx.Params(":id")
	groups,err := auth.GetGroupChildren(groupId)
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

//For Get POST Info Only
func GetGroupPostInfo(ctx *macaron.Context)(args map[string]interface{}, err error){
	body, err := utils.GetPostBody(ctx)
	tokenStr := ctx.Req.Header.Get("token")
	user, err := auth.GetUserByTokenString(tokenStr)
	if err != nil{
		return
	}
	args = make(map[string]interface{},0)
	for k,v := range body {
		switch k {
		case "id":
			args["id"] = v.(string)
		case "name":
			args["name"] = v.(string)
		case "description":
			args["description"] = v.(string)
		case "admin_user_id":
			args["admin_user_id"] = v.(string)
		case "parent_id":
			args["parent_id"] = v.(string)
		case "status":
			args["status"] = v.(string)
		case "end_time":
			args["end_time"] = utils.InterfaceToInt(v)
		}
	}
	args["current_user"] = user.Name
	return
}