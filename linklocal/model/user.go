package model

import(
	"time"
	"fmt"
	"github.com/surma-dump/gouuid"
	"linklocal/utils"
)

type User struct{
	Id	string	`xorm:"pk"`
	Name	string	`xorm:"unique"`
	TrueName string
	Password string
	Type	string
	Email	string	`xorm:"unique"`
	Tel	string	`xorm:"varchar(30)"`
	Language string	`xorm:"varchar(20)"`
	Description	string
	EndTime	time.Time
	EndAction	string	`xorm:"varchar(20)"`
	CreatedBy    string `xorm:"varchar(50)"`
	CreatedDate  time.Time
	ModifiedBy   string `xorm:"varchar(50)"`
	ModifiedDate time.Time
	Groups []Group
}

func GetUserByID(id string)(user *User,err error){
	users := make([]User, 0)
	err = x.Where("id=?", id).Find(&users)
	if err != nil {
		return
	}
	if len(users) == 0 {
		err = fmt.Errorf("User Not Exist")
		return
	}
	user = &users[0]
	return
}

func GetUserByName(name string)(user *User,err error){
	users := make([]User, 0)
	err = x.Where("name=?", name).Find(&users)
	if err != nil {
		return
	}
	if len(users) == 0 {
		err = fmt.Errorf("User Not Exist")
		return
	}
	user = &users[0]
	return
}

func GetUserByEmail(email string)(user *User,err error){
	users := make([]User, 0)
	err = x.Where("email=?", email).Find(&users)
	if err != nil {
		return
	}
	if len(users) == 0 {
		err = fmt.Errorf("User Not Exist")
		return
	}
	user = &users[0]
	return
}

func ListUsers()(users []User,err error){
	users = make([]User,0)
	err = x.Find(&users)
	if err!=nil{
		return
	}
	return
}

func CreateUser(user *User)(err error){
	_,err = GetUserByName(user.Name)
	if err == nil{
		err = fmt.Errorf("User Name Conflict")
	}
	_,err = GetUserByEmail(user.Email)
	if err == nil{
		err = fmt.Errorf("User Email Conflict")
	}
	uuid := gouuid.New().String()
	for {
		_, err = GetUserByID(uuid)
		if err == nil {
			uuid = gouuid.New().String()
		} else {
			break
		}
	}
	user.Id = uuid

	affect, err := x.InsertOne(user)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Insert User Faild")
		return
	}
	return
}

func DeleteUser(user *User) (err error) {
	affect, err := x.Id(user.Id).Delete(user)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Delete User Faild")
		return
	}
	return
}

func UpdateUser(user *User) (err error) {
	affect, err := x.Id(user.Id).Update(user)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Update User Faild")
		return
	}
	return
}

func ChangeUserPassword(id,password string)(err error){
	user,err := GetUserByID(id)
	if err != nil{
		return
	}
	user.Password = utils.GetMd5s(password)
	UpdateUser(user)
	return
}

func ChangeUserEndTime(id string,endTime int)(err error){
	user,err := GetUserByID(id)
	if err != nil{
		return
	}
	endTimeDuration := Eternal
	if endTime != -1{
		endTimeDuration = time.Duration(24*endTime)*time.Hour
	}
	user.EndTime = time.Now().Add(endTimeDuration)
	UpdateUser(user)
	return
}
