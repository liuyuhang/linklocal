package controller

import(
	"gopkg.in/macaron.v1"
	"linklocal/utils"
	"linklocal/auth"
	"fmt"
)

func Login(ctx *macaron.Context){
	loginInfo := make(map[string]interface{})
	args,err := GetLoginInfo(ctx)
	if err != nil{
		ctx.JSON(500,err.Error())
		return
	}
	_,ok := args["username"]
	_,ok1 := args["password"]
	if !(ok && ok1){
		err := fmt.Errorf("username and password should be set")
		loginInfo["error"] = err.Error()
		loginInfo["success"] = false
		loginInfo["user"] = args["username"].(string)
		loginInfo["token"] = ""
		ctx.JSON(401, loginInfo)
		return
	}

	token,err := auth.LoginAuthorization(args["username"].(string),args["password"].(string))
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


func GetLoginInfo(ctx *macaron.Context) (args map[string]interface{}, err error) {
	args = make(map[string]interface{},0)
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
		}
	}
	return
}
