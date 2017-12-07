package model

import(
	"time"
	"fmt"
	"github.com/surma-dump/gouuid"
)

type Group struct{
	Id	string	`xorm:"pk"`
	Name	string	`xorm:"unique"`
	Description	string
	AdminUserId string
	AdminUserName string
	ParentId	string
	ParentName string
	Status	string
	EndTime	time.Time
	EndAction	string	`xorm:"varchar(20)"`
	CreatedBy    string `xorm:"varchar(50)"`
	CreatedDate  time.Time
	ModifiedBy   string `xorm:"varchar(50)"`
	ModifiedDate time.Time
}

func GetGroupByID(id string)(group *Group,err error){
	groups := make([]Group, 0)
	err = x.Where("id=?", id).Find(&groups)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		err = fmt.Errorf("Group Not Exist")
		return
	}
	group = &groups[0]
	if group.AdminUserId != ""{
		user,err := GetUserByID(group.AdminUserId)
		if err!=nil{
			group.AdminUserName = "NotFound"
		}else{
			group.AdminUserName = user.Name
		}
	}
	if group.ParentId != "" {
		gp,err := GetGroupByID(group.ParentId)
		if err !=nil{
			group.ParentName = ""
		}else{
			group.ParentName = gp.Name

		}
	}
	return
}

func GetGroupByName(name string)(group *Group,err error){
	groups := make([]Group, 0)
	err = x.Where("name=?", name).Find(&groups)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		err = fmt.Errorf("Group Not Exist")
		return
	}
	group = &groups[0]
	if group.AdminUserId != ""{
		user,err := GetUserByID(group.AdminUserId)
		if err!=nil{
			group.AdminUserName = "NotFound"
		}else{
			group.AdminUserName = user.Name
		}
	}
	if group.ParentId != "" {
		gp,err := GetGroupByID(group.ParentId)
		if err !=nil{
			group.ParentName = ""
		}else{
			group.ParentName = gp.Name
		}
	}
	return
}

func ListGroups()(groups []Group,err error){
	groups = make([]Group,0)
	err = x.Find(&groups)
	if err!=nil{
		return
	}
	for k,_:= range groups{
		if groups[k].AdminUserId != ""{
			user,err := GetUserByID(groups[k].AdminUserId)
			if err!=nil{
				groups[k].AdminUserName = "NotFound"
			}else{
				groups[k].AdminUserName = user.Name
			}
		}
		if groups[k].ParentId != "" {
			gp,err := GetGroupByID(groups[k].ParentId)
			if err !=nil{
				groups[k].ParentName = ""
			}else{
				groups[k].ParentName = gp.Name
			}
		}
	}
	return
}

func CreateGroup(group *Group)(err error){
	_,err = GetGroupByName(group.Name)
	if err == nil{
		err = fmt.Errorf("Group Name Conflict")
	}
	uuid := gouuid.New().String()
	for {
		_, err = GetGroupByID(uuid)
		if err == nil {
			uuid = gouuid.New().String()
		} else {
			break
		}
	}
	group.Id = uuid

	affect, err := x.InsertOne(group)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Insert Group Faild")
		return
	}
	return
}

func DeleteGroup(group *Group) (err error) {
	affect, err := x.Id(group.Id).Delete(group)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Delete Group Faild")
		return
	}
	return
}

func UpdateGroup(group *Group) (err error) {
	affect, err := x.Id(group.Id).Update(group)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Update Group Faild")
		return
	}
	return
}

func GetGroupParent(group *Group)(pgroup *Group,err error){
	parentID := group.ParentId
	pgroup,err = GetGroupByID(parentID)
	return
}

func GetGroupChildren(group *Group)(children []Group,err error){
	parentID := group.Id
	children = make([]Group, 0)
	err = x.Where("parent_id=?", parentID).Find(&children)
	if err != nil {
		return
	}
	if len(children) == 0 {
		err = fmt.Errorf("Group Not Exist")
		return
	}
	return
}

func ChangeGroupEndTime(id string,endTime int)(err error){
	group,err := GetGroupByID(id)
	if err != nil{
		return
	}
	endTimeDuration := Eternal
	if endTime != -1{
		endTimeDuration = time.Duration(24*endTime)*time.Hour
	}
	group.EndTime = time.Now().Add(endTimeDuration)
	UpdateGroup(group)
	return
}