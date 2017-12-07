package model

import(
	"github.com/ccwings/log"
	"github.com/go-xorm/xorm"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var(
	x *xorm.Engine
	Eternal = 24*365*100*time.Hour
)

func InitModelDB(dataDriver,dataSource string)(err error){
	log.Info("Init ModelDB ...")
	x,err = xorm.NewEngine(dataDriver,dataSource)
	if err!=nil{
		log.Error(err.Error())
	}
	for{
		err = x.Ping()
		if err!=nil{
			time.Sleep(5*time.Second)
		}else{
			break
		}
	}
	defer initModelData()
	err = x.Sync(new(User),new(Group),new(UserGroup))
	if err!= nil{
		log.Error(err.Error())
		return
	}
	return
}

