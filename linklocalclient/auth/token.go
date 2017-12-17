package auth

import (
	"time"
	"linklocal/utils"
	"fmt"
	"linklocal/model"
	"github.com/ccwings/log"
)

var (
	tokenList=make(map[string]*Token)
)

const expiryDelta = 1800 * time.Second

type Token struct{
	AccessToken string `json:"access_token"`
	TokenType	string	`json:"token_type"`
	Expiry	time.Time	`json:"expiry"`
	User	string	`json:"user"`
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

func CreateToken(username string)(t *Token){
	t = new(Token)
	t.AccessToken = utils.RandSeq(50)
	t.TokenType = "Brarer"
	t.Expiry = time.Now().Add(expiryDelta)
	t.User = username
	return t
}

func GetTokenByTokenString(token_str string)(token *Token,err error){
	for _,v := range(tokenList){
		if v.AccessToken == token_str{
			return v,nil
		}
	}
	err = fmt.Errorf("Token Not Found")
	return nil,err
}

func GetUserByTokenString(token_str string) (user *model.User,err error){
	token,err := GetTokenByTokenString(token_str)
	if err!=nil{
		log.Error(err)
		return
	}
	user,err = model.GetUserByName(token.User)
	if err!=nil{
		log.Error(err)
		return
	}
	return
}

// expired reports whether the token is expired.
func (t *Token) expired() bool {
	if t==nil{
		return false
	}
	if t.Expiry.IsZero() {
		return false
	}
	//return t.Expiry.Add(-expiryDelta).Before(time.Now())
	if t.Expiry.Before(time.Now()){
		t.Expiry = time.Now()
		return true
	}
	return false
}

func TokenAuthorization(token string) bool{

	for k,v := range tokenList{
		log.Debug(k,v)
	}

	t,ok := tokenList[token]
	if !ok{
		return false
	}else{
		if t.Valid(){
			log.Debug("Valid sucess")
			return true
		}else{
			log.Debug("Valid faild")

			delete(tokenList,token)
			return false
		}
	}
	return false
}

func InitDevToken(){
	t := CreateToken("admin")
	t.AccessToken = "123456"
	tokenList["123456"] = t
	return
}