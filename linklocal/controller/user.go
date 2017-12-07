package controller

import (
	"gopkg.in/macaron.v1"
	"linklocal/auth"
	"linklocal/utils"
	"linklocal/model"
)

func CreateUser(ctx *macaron.Context){
	info := make(map[string]interface{})
	args,err := GetUserPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["user"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	user, err := auth.CreateUser(args)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["user"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["user"] = user
	ctx.JSON(200,info)
	return
}

func DeleteUser(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	err := auth.DeleteUser(userId)
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

func UpdateUser(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	args,err := GetUserPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["user"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	user,err := auth.UpdateUser(userId,args)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["user"] = make(map[string]interface{},0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["user"] = user
	ctx.JSON(200,info)
	return

}

func ListUsers(ctx *macaron.Context){
	info := make(map[string]interface{})
	users,err := auth.ListUsers()
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["users"] = make([]model.User,0)
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["users"] = users
	ctx.JSON(200,info)
	return
}

func GetUser(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	user,err := auth.GetUser(userId)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		info["user"] = make(map[string]interface{})
		ctx.JSON(500,info)
		return
	}
	info["error"] = nil
	info["success"] = true
	info["user"] = user
	ctx.JSON(200,info)
	return
}

func ChangeUserPassword(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	args,err := GetUserPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.ChangeUserPassword(userId,args)
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

func ChangeUserEndTime(ctx *macaron.Context){
	info := make(map[string]interface{})
	userId := ctx.Params(":id")
	args,err := GetUserPostInfo(ctx)
	if err!=nil{
		info["error"] = err.Error()
		info["success"] = false
		ctx.JSON(500,info)
		return
	}
	err = auth.ChangeUserEndTime(userId,args)
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
func GetUserPostInfo(ctx *macaron.Context)(args map[string]interface{}, err error){
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
		case "password":
			args["password"] = v.(string)
		case "email":
			args["email"] = v.(string)
		case "description":
			args["description"] = v.(string)
		case "tel":
			args["tel"] = v.(string)
		case "language":
			args["language"] = v.(string)
		case "type":
			args["type"] = v.(string)
		case "end_time":
			args["end_time"] = utils.InterfaceToInt(v)
		case "true_name":
			args["true_name"] = v.(string)

		}
	}
	args["current_user"] = user.Name
	return
}
