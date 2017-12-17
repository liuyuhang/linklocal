package controller

import (
	"gopkg.in/macaron.v1"
	"linklocal/utils"
	"linklocal/auth"
	"fmt"
)

func Login(ctx *macaron.Context) {
	loginInfo := make(map[string]interface{})
	args, err := GetLoginInfo(ctx)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	_, ok := args["username"]
	_, ok1 := args["password"]
	if !(ok && ok1) {
		err := fmt.Errorf("username and password should be set")
		loginInfo["error"] = err.Error()
		loginInfo["success"] = false
		loginInfo["user"] = args["username"].(string)
		loginInfo["token"] = ""
		ctx.JSON(401, loginInfo)
		return
	}

	token, err := auth.LoginAuthorization(args["username"].(string), args["password"].(string))
	if err != nil {
		loginInfo["error"] = err.Error()
		loginInfo["success"] = false
		loginInfo["user"] = args["username"].(string)
		loginInfo["token"] = ""
		ctx.JSON(401, loginInfo)
		return
	}
	loginInfo["error"] = nil
	loginInfo["success"] = true
	loginInfo["user"] = args["username"].(string)
	loginInfo["token"] = token
	ctx.JSON(200, loginInfo)
	return
}

func Logout(ctx *macaron.Context) {
	logoutInfo := make(map[string]interface{})
	token := ctx.Req.Header.Get("token")
	err := auth.LogoutAuthorization(token)
	if err != nil {
		logoutInfo["error"] = err.Error()
		logoutInfo["success"] = false
		ctx.JSON(401, logoutInfo)
		return
	}
	logoutInfo["error"] = nil
	logoutInfo["success"] = true
	ctx.JSON(200, logoutInfo)
	return
}

func Register(ctx *macaron.Context) {
	registerInfo := make(map[string]interface{})
	args, err := GetLoginInfo(ctx)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	_, ok := args["username"]
	_, ok1 := args["password"]
	_, ok2 := args["email"]
	if !(ok && ok1 && ok2) {
		err := fmt.Errorf("username, email and password should be set")
		registerInfo["error"] = err.Error()
		registerInfo["success"] = false
		registerInfo["user"] = args["username"].(string)
		registerInfo["token"] = ""
		ctx.JSON(500, registerInfo)
		return
	}

	err = auth.Register(args["username"].(string), args["email"].(string), args["password"].(string))
	if err != nil{
		registerInfo["error"] = err.Error()
		registerInfo["success"] = false
		registerInfo["user"] = args["username"].(string)
		registerInfo["token"] = ""
		ctx.JSON(500, registerInfo)
		return
	}

	token, err := auth.LoginAuthorization(args["username"].(string), args["password"].(string))
	if err != nil {
		registerInfo["error"] = err.Error()
		registerInfo["success"] = false
		registerInfo["user"] = args["username"].(string)
		registerInfo["token"] = ""
		ctx.JSON(401, registerInfo)
		return
	}
	registerInfo["error"] = nil
	registerInfo["success"] = true
	registerInfo["user"] = args["username"].(string)
	registerInfo["token"] = token
	ctx.JSON(200, registerInfo)
	return
}


func CheckUsername(ctx *macaron.Context) {
	checkInfo := make(map[string]interface{})
	args, err := GetLoginInfo(ctx)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	_, ok := args["username"]
	if !(ok) {
		err := fmt.Errorf("username should be set")
		checkInfo["error"] = err.Error()
		checkInfo["success"] = false
		ctx.JSON(500, checkInfo)
		return
	}
	err = auth.CheckUsername(args["username"].(string))
	if err != nil {
		checkInfo["error"] = err.Error()
		checkInfo["success"] = false
		ctx.JSON(500, checkInfo)
		return
	}
	checkInfo["error"] = nil
	checkInfo["success"] = true
	ctx.JSON(200, checkInfo)
	return
}
func CheckEmail(ctx *macaron.Context) {
	checkInfo := make(map[string]interface{})
	args, err := GetLoginInfo(ctx)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	_, ok := args["email"]
	if !(ok) {
		err := fmt.Errorf("email should be set")
		checkInfo["error"] = err.Error()
		checkInfo["success"] = false
		ctx.JSON(500, checkInfo)
		return
	}
	err = auth.CheckEmail(args["email"].(string))
	if err != nil {
		checkInfo["error"] = err.Error()
		checkInfo["success"] = false
		ctx.JSON(500, checkInfo)
		return
	}
	checkInfo["error"] = nil
	checkInfo["success"] = true
	ctx.JSON(200, checkInfo)
	return
}

func GetLoginInfo(ctx *macaron.Context) (args map[string]interface{}, err error) {
	args = make(map[string]interface{}, 0)
	body, err := utils.GetPostBody(ctx)
	if err != nil {
		return
	}
	for k, v := range body {
		switch k {
		case "username":
			args["username"] = v.(string)
		case "password":
			args["password"] = v.(string)
		case "domainid":
			args["domain"] = v.(string)
		case "email":
			args["email"] = v.(string)
		}
	}
	return
}
