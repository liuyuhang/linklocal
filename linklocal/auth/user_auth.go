package auth

import (
	"linklocal/model"
	"strings"
	"fmt"
	"github.com/ccwings/log"
	"github.com/ccwings/conf"
	"linklocal/utils"
	"time"
)

var(
	auth_type = ""
)

func LoginAuthorization(username,password string)(token *Token,err error){
	log.Debug("in login auth")
	GetLoginType()
	if auth_type == "local"{
		err = localAuthorization(username,password)
		if err!=nil{
			log.Debug(err)
			return
		}
	}else{
		err = fmt.Errorf("Other Login Type Not Support.")
		return
	}
	token = CreateToken(username)
	_,ok := tokenList[token.AccessToken]
	if !ok{
		tokenList[token.AccessToken] = token
	}

	return tokenList[token.AccessToken] ,nil
}

func LogoutAuthorization(token string)error{
	log.Debug("in logout auth")
	delete(tokenList,token)
	return nil
}

func localAuthorization(username,password string)(err error){
	var user *model.User
	if strings.Contains(username,"@"){
		user,err = model.GetUserByEmail(username)
		if err!=nil{
			return err
		}
	}else{
		user,err = model.GetUserByName(username)
		if err!=nil{
			return err
		}
	}

	if user.Password != utils.GetMd5s(password){
		return fmt.Errorf("Authorization Faild")
	}
	if user.EndTime.Before(time.Now()){
		return fmt.Errorf("User Not Active")
	}
	return nil
}

func GetLoginType()(err error){
	if auth_type != ""{
		return nil
	}
	config, err := conf.ReadConfigFile("config.conf")
	if err != nil {
		config, err = conf.ReadConfigFile("/etc/ccwings/bigbang/config.conf")
		if err != nil {
			log.Debug(err)
		}
	}
	auth_type,err = config.GetString("default","auth_type")
	if err!=nil{
		auth_type = "local"
	}
	return
}