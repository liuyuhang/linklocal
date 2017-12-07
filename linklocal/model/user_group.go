package model

import (
	"time"
	"fmt"
	"github.com/surma-dump/gouuid"
)

type UserGroup struct{
	Id	string	`xorm:"pk"`
	UserId	string
	GroupId	string
	Type	string
	Priority string
	EndTime	time.Time
	EndAction	string	`xorm:"varchar(20)"`
	CreatedBy    string `xorm:"varchar(50)"`
	CreatedDate  time.Time
}


func GetUserGroupByID(id string)(userGroup *UserGroup,err error){
	userGroups := make([]UserGroup, 0)
	err = x.Where("id=?", id).Find(&userGroups)
	if err != nil {
		return
	}
	if len(userGroups) == 0 {
		err = fmt.Errorf("UserGroup Not Exist")
		return
	}
	userGroup = &userGroups[0]
	return
}

func GetUserGroupByUG(userID,groupID string)(userGroup *UserGroup,err error){
	userGroups := make([]UserGroup, 0)
	err = x.Where("user_id=? and group_id=?", userID,groupID).Find(&userGroups)
	if err != nil {
		return
	}
	if len(userGroups) == 0 {
		err = fmt.Errorf("UserGroup Not Exist")
		return
	}
	userGroup = &userGroups[0]
	return
}

func CreateUserGroup(userGroup *UserGroup)(err error){
	_,err = GetUserGroupByUG(userGroup.UserId,userGroup.GroupId)
	if err == nil{
		err = fmt.Errorf("UserGroup Conflict")
	}
	uuid := gouuid.New().String()
	for {
		_, err = GetUserGroupByID(uuid)
		if err == nil {
			uuid = gouuid.New().String()
		} else {
			break
		}
	}
	userGroup.Id = uuid

	affect, err := x.InsertOne(userGroup)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Insert UserGroup Faild")
		return
	}
	return
}

func DeleteUserGroup(userGroup *UserGroup) (err error) {
	affect, err := x.Id(userGroup.Id).Delete(userGroup)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Delete UserGroup Faild")
		return
	}
	return
}

func UpdateUserGroup(userGroup *UserGroup) (err error){
	affect, err := x.Id(userGroup.Id).Update(userGroup)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Update UserGroup Faild")
		return
	}
	return
}

func GetUserGroupByUserID(userID string)(userGroups []UserGroup,err error){
	userGroups = make([]UserGroup, 0)
	err = x.Where("user_id=?", userID).Find(&userGroups)
	if err != nil {
		return
	}
	if len(userGroups) == 0 {
		err = fmt.Errorf("UserGroup Not Exist")
		return
	}
	return
}

func GetUserGroupByGroupID(groupID string)(userGroups []UserGroup,err error){
	userGroups = make([]UserGroup, 0)
	err = x.Where("group_id=?", groupID).Find(&userGroups)
	if err != nil {
		return
	}
	if len(userGroups) == 0 {
		err = fmt.Errorf("UserGroup Not Exist")
		return
	}
	return
}
