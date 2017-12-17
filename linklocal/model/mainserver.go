package model

import(
	"time"
	"fmt"
	"github.com/surma-dump/gouuid"
)


type MainServer struct{
	Id	string	`xorm:"pk"`
	Name	string	`xorm:"unique"`
	Description	string
	LastCheckTime time.Time
	CreatedBy    string `xorm:"varchar(50)"`
	CreatedDate  time.Time
	ModifiedBy   string `xorm:"varchar(50)"`
	ModifiedDate time.Time

}

func GetMainServerByName(name string)(mainServer *MainServer,err error){
	mainServers := make([]MainServer, 0)
	err = x.Where("name=?", name).Find(&mainServers)
	if err != nil {
		return
	}
	if len(mainServers) == 0 {
		err = fmt.Errorf("MainServer Not Exist")
		return
	}
	mainServer = &mainServers[0]
	return
}

func GetMainServerById(id string)(mainServer *MainServer,err error){
	mainServers := make([]MainServer, 0)
	err = x.Where("id=?", id).Find(&mainServers)
	if err != nil {
		return
	}
	if len(mainServers) == 0 {
		err = fmt.Errorf("MainServer Not Exist")
		return
	}
	mainServer = &mainServers[0]
	return
}

func ListMainServer()(mainServers []MainServer,err error){
	mainServers = make([]MainServer,0)
	err = x.Find(&mainServers)
	if err!=nil{
		return
	}
	return
}

func CreateMainServer(mainServer *MainServer)(err error){
	_,err = GetMainServerByName(mainServer.Name)
	if err == nil{
		err = fmt.Errorf("MainServer Name Conflict")
	}
	uuid := gouuid.New().String()
	for {
		_, err = GetMainServerById(uuid)
		if err == nil {
			uuid = gouuid.New().String()
		} else {
			break
		}
	}
	mainServer.Id = uuid

	affect, err := x.InsertOne(mainServer)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Insert MainServer Faild")
		return
	}
	return
}

func UpdateMainServer(mainServer *MainServer) (err error) {
	affect, err := x.Id(mainServer.Id).Update(mainServer)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Update MainServer Faild")
		return
	}
	return
}

func DeleteMainServer(mainServer *MainServer) (err error) {
	affect, err := x.Id(mainServer.Id).Delete(mainServer)
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Delete MainServer Faild")
		return
	}
	return
}