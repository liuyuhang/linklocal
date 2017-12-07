package auth

import (
	"linklocal/model"
	"time"
	"linklocal/utils"
	"github.com/ccwings/log"
	"fmt"
)

func CreateGroup(args map[string]interface{})(group *model.Group,err error){
	group = new(model.Group)
	_,ok := args["admin_user_id"]
	if args["admin_user_id"] == "" || !ok{
		user,err := model.GetUserByName(args["current_user"].(string))
		if err!=nil{
			log.Error(err.Error())
		}
		args["admin_user_id"] = user.Id
	}
	for k,v := range args{
		switch k {
		case "name":
			group.Name = v.(string)
		case "admin_user_id":
			group.AdminUserId = v.(string)
		case "parent_id":
			group.ParentId = v.(string)
		case "status":
			group.Status = v.(string)
		case "description":
			group.Description = v.(string)
		case "end_time":
			group.EndTime = time.Now().Add(time.Duration(utils.InterfaceToInt(v)*24)*time.Hour)
		}
	}
	_,err = model.GetGroupByID(group.ParentId)
	if err !=nil{
		err = fmt.Errorf("Parent Group Not Found.")
		return
	}

	group.CreatedBy = args["current_user"].(string)
	group.CreatedDate = time.Now()
	group.EndAction = "keep"
	err = model.CreateGroup(group)
	if err != nil{
		log.Error(err.Error())
	}
	return
}

func DeleteGroup(groupID string)(err error){
	// todo Check if group has users?
	log.Debug(groupID)
	group,err := model.GetGroupByID(groupID)
	if err!=nil{
		err = fmt.Errorf("Group Not Found")
		return
		}

	groupChild,err := model.GetGroupChildren(group)
	if err == nil || len(groupChild)>0{
		err= fmt.Errorf("Group Has At Least One Child Group. Can Not Be Deleted.")
		return
		}

	err = model.DeleteGroup(group)
	return
}

func UpdateGroup(groupID string,args map[string]interface{})(group *model.Group,err error){
	group,err = model.GetGroupByID(groupID)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	for k,v := range args{
		switch k {
		case "name":
			group.Name = v.(string)
		case "admin_user_id":
			group.AdminUserId = v.(string)

		case "parent_id":
			group.ParentId = v.(string)
			pg,err := model.GetGroupByID(group.ParentId)
			if err!= nil{
				err = fmt.Errorf("Parent Group Not Found")
				return group,err
			}
			group.ParentName = pg.Name
		case "status":
			group.Status = v.(string)
		case "description":
			group.Description = v.(string)
		}
	}
	group.ModifiedBy = args["current_user"].(string)
	group.ModifiedDate = time.Now()
	err = model.UpdateGroup(group)
	if err != nil{
		log.Error(err.Error())
	}
	return
}

func ListGroups()(groups []model.Group,err error){
	groups,err = model.ListGroups()
	return
}

func GetGroup(id string)(group *model.Group,err error){
	group,err = model.GetGroupByID(id)
	return
}

func GetGroupParent(id string)(pgroup *model.Group,err error){
	group,err := model.GetGroupByID(id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	pgroup,err = model.GetGroupParent(group)
	if err != nil{
		log.Error(err.Error())
	}
	return
}

func GetGroupChildren(id string)(groupChildren []model.Group,err error){
	group,err := model.GetGroupByID(id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	groupChildren,err = model.GetGroupChildren(group)
	if err != nil{
		log.Error(err.Error())
	}
	return
}


func ChangeGroupEndTime(id string, args map[string]interface{})(err error){
	endTime := 0
	for k,v := range args{
		switch k {
		case "end_time":
			endTime = utils.InterfaceToInt(v)
		}
	}
	if endTime == 0{
		err = fmt.Errorf("Args EndTime Not Set.")
		return
	}
	model.ChangeGroupEndTime(id,endTime)
	return
}