package auth

import(
	"linklocal/model"
	"time"
	"linklocal/utils"
	"github.com/ccwings/log"
	"fmt"
)

func CreateUser(args map[string]interface{})(user *model.User,err error){
	user = new(model.User)
	for k,v := range args{
		switch k {
		case "name":
			user.Name = v.(string)
		case "password":
			user.Password = utils.GetMd5s(v.(string))
		case "type":
			user.Type = v.(string)
		case "email":
			user.Email = v.(string)
		case "description":
			user.Description = v.(string)
		case "end_time":
			user.EndTime = time.Now().Add(time.Duration(utils.InterfaceToInt(v)*24)*time.Hour)
		case "tel":
			user.Tel = v.(string)
		case "true_name":
			user.TrueName = v.(string)
		}
	}
	user.CreatedBy = args["current_user"].(string)
	user.CreatedDate = time.Now()
	user.EndAction = "keep"
	err = model.CreateUser(user)
	if err != nil{
		log.Error(err.Error())
	}
	user.Groups = make([]model.Group,0)
	user.Password = "******"
	return
}

func DeleteUser(user_id string)(err error){
	// todo Check if user has Admin To Group
	log.Debug(user_id)
	user,err := model.GetUserByID(user_id)
	if err!=nil{
		err = fmt.Errorf("User Not Found")
		return
	}
	err = model.DeleteUser(user)
	return
}

func UpdateUser(user_id string,args map[string]interface{})(user *model.User,err error){
	user,err = model.GetUserByID(user_id)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	for k,v := range args{
		switch k {
		case "name":
			user.Name = v.(string)
		case "type":
			user.Type = v.(string)
		case "email":
			user.Email = v.(string)
		case "description":
			user.Description = v.(string)
		case "tel":
			user.Tel = v.(string)
		case "true_name":
			user.TrueName = v.(string)
		}
	}
	user.ModifiedBy = args["current_user"].(string)
	user.ModifiedDate = time.Now()
	err = model.UpdateUser(user)
	if err != nil{
		log.Error(err.Error())
	}
	user.Password = "******"
	user.Groups,err = GetUserBindingGroups(user.Id)
	if err !=nil{
		user.Groups = make([]model.Group,0)
	}
	return user,nil
}

func ListUsers()(users []model.User,err error){
	users,err = model.ListUsers()
	for k,_ := range users{
		users[k].Password = "******"
		ug,err1 := GetUserBindingGroups(users[k].Id)
		if err1 !=nil{
			users[k].Groups = make([]model.Group,0)
			continue
		}
		users[k].Groups = ug
	}
	return
}

func GetUser(id string)(user *model.User,err error){
	user,err = model.GetUserByID(id)
	user.Password = "******"
	ug,err1 := GetUserBindingGroups(user.Id)
	if err1 !=nil{
		user.Groups = make([]model.Group,0)
		return
	}
	user.Groups = ug
	return
}

func ChangeUserPassword(id string, args map[string]interface{})(err error){
	password := ""
	for k,v := range args{
		switch k {
		case "password":
			password = v.(string)
		}
	}
	if password == ""{
		err = fmt.Errorf("Args Password Not Set.")
		return
	}
	model.ChangeUserPassword(id,password)
	return
}

func ChangeUserEndTime(id string, args map[string]interface{})(err error){
	endtime := 0
	for k,v := range args{
		switch k {
		case "end_time":
			endtime = utils.InterfaceToInt(v)
		}
	}
	if endtime == 0{
		err = fmt.Errorf("Args EndTime Not Set.")
		return
	}
	model.ChangeUserEndTime(id,endtime)
	return
}