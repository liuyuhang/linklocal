package model

import (
	"linklocal/utils"
	"time"
	"github.com/ccwings/log"
)

func initModelData(){
	initUser()
	initGroup()
	initUserGroup()
	return
}

func initUser(){
	user1 := new(User)
	user1.Name = "admin"
	user1.Password = utils.GetMd5s("123")
	user1.Type = "admin"
	user1.TrueName = "SuperAdmin"
	user1.Email = "admin@linklocal.cn"
	user1.Description = "SuperAdmin"
	user1.CreatedBy = "OS"
	user1.CreatedDate = time.Now()
	user1.EndTime =time.Now().Add(Eternal)
	user1.EndAction = "keep"
	user1.Tel = "1234567890"
	user1.Language = "cn"
	err := CreateUser(user1)
	if err != nil{
		log.Error(err.Error())
	}
	return
}

func initGroup(){
	adminUser, err := GetUserByName("admin")
	if err!=nil{
		log.Error("Admin User Not Found")
		return
	}
	adminGroup := new(Group)
	adminGroup.Name = "AdminGroup"
	adminGroup.Status = "active"
	adminGroup.ParentId = ""
	adminGroup.AdminUserId = adminUser.Id
	adminGroup.CreatedBy = "OS"
	adminGroup.CreatedDate = time.Now()
	adminGroup.EndAction = "keep"
	adminGroup.EndTime = time.Now().Add(Eternal)
	adminGroup.Description = "AdminGroup"
	err = CreateGroup(adminGroup)
	if err != nil{
		log.Error(err.Error())
	}
	return
}

func initUserGroup(){
	adminUser, err := GetUserByName("admin")
	if err!=nil{
		log.Error("Admin User Not Found")
		return
	}
	adminGroup,err := GetGroupByName("AdminGroup")
	if err!=nil{
		log.Error("Admin User Group Not Found")
		return
	}
	_,err = GetUserGroupByUG(adminUser.Id,adminGroup.Id)
	if err== nil{
		log.Info("UserGroup exist")
		return
	}
	adminUserGroup := new(UserGroup)
	adminUserGroup.UserId = adminUser.Id
	adminUserGroup.GroupId = adminGroup.Id
	adminUserGroup.EndTime=time.Now().Add(Eternal)
	adminUserGroup.EndAction = "keep"
	adminUserGroup.CreatedDate = time.Now()
	adminUserGroup.CreatedBy = "OS"
	adminUserGroup.Type = "admin"
	adminUserGroup.Priority = "1"
	err = CreateUserGroup(adminUserGroup)
	if err != nil{
		log.Error(err.Error())
	}
	return
}