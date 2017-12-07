package auth

import (
	"linklocal/model"
	"github.com/ccwings/log"
	"time"
	"linklocal/utils"
)

func GetGroupBindingUsers(id string)(users []model.User,err error){
	userGroups,err := model.GetUserGroupByGroupID(id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	users = make([]model.User,0)
	for _,v := range userGroups{
		user,err := model.GetUserByID(v.UserId)
		if err!=nil{
			continue
		}
		users = append(users,*user)
	}

	return
}

func GetUserBindingGroups(id string)(groups []model.Group,err error){
	userGroups,err := model.GetUserGroupByUserID(id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	groups = make([]model.Group,0)
	for _,v := range userGroups{
		group,err := model.GetGroupByID(v.GroupId)
		if err!=nil{
			continue
		}
		groups = append(groups,*group)
	}

	return
}

func BindUserToGroup(groupId string,args map[string]interface{})(err error) {
	userGroup := new(model.UserGroup)
	for k,v := range args{
		switch k {
		case "user_id":
			userGroup.UserId = v.(string)
		case "type":
			userGroup.Type = v.(string)
		case "priority":
			userGroup.Priority = v.(string)
		case "end_time":
			userGroup.EndTime = time.Now().Add(time.Duration(utils.InterfaceToInt(v)*24)*time.Hour)
		}
	}
	userGroup.GroupId = groupId
	userGroup.CreatedBy = args["current_user"].(string)
	userGroup.CreatedDate = time.Now()
	err = model.CreateUserGroup(userGroup)
	if err != nil{
		log.Error(err.Error())
	}
	return

}

func UnbindUserToGroup(groupId,userId string)(err error) {
	userGroup,err := model.GetUserGroupByUG(userId,groupId)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	err = model.DeleteUserGroup(userGroup)
	if err != nil{
		log.Error(err.Error())
	}
	return

}

func SetUserTypeForGroup(groupId,userId,userType string)(err error){
	userGroup,err := model.GetUserGroupByUG(userId,groupId)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	userGroup.Type = userType
	err = model.UpdateUserGroup(userGroup)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	return nil
}

func SetUserPriorityForGroup(groupId,userId,userPriority string)(err error){
	userGroup,err := model.GetUserGroupByUG(userId,groupId)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	userGroup.Priority = userPriority
	err = model.UpdateUserGroup(userGroup)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	return nil
}