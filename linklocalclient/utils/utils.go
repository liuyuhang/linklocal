package utils

import(
	"math/rand"
	"time"
	"reflect"
	"strings"
	"strconv"
	"github.com/ccwings/conf"
	"github.com/ccwings/log"
)

// rand a n-length string
func RandSeq(n int)string{
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b :=make([]rune,n)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i:= range b{
		b[i] = letters[r1.Intn(len(letters))]
	}
	return string(b)
}


func InterfaceToBool(v interface{})(bool){
	if reflect.TypeOf(v).String()==reflect.String.String(){
		if strings.ToLower(v.(string)) == "true"{
			return true
		}
	}
	if reflect.TypeOf(v).String() == reflect.Bool.String(){
		return v.(bool)
	}
	return false
}

func InterfaceToInt(v interface{})(int){
	if reflect.TypeOf(v).String()==reflect.Float64.String(){
		return int(v.(float64))
	}
	if reflect.TypeOf(v).String() == reflect.Int.String(){
		return v.(int)
	}
	if reflect.TypeOf(v).String() == reflect.String.String(){
		i,err := strconv.Atoi(v.(string))
		if err!=nil{
			return 0
		}
		return i
	}
	return 0
}

func BoolToString(v bool)(string){
	if v{
		return "True"
	}else{
		return "False"
	}
	return "False"
}

func StringToBool(v string)(bool){
	if strings.ToLower(v)=="true"{
		return true
	}
	return false
}

// count ips for format: start:192.168.0.10  end:192.168.0.20
func CountIPs(startS,endS string)(count int, err error){
	ipStart := strings.Split(startS,".")
	ipEnd := strings.Split(endS,".")
	ipStart0, _ := strconv.Atoi(ipStart[0])
	ipStart1, _ := strconv.Atoi(ipStart[1])
	ipStart2, _ := strconv.Atoi(ipStart[2])
	ipStart3, _ := strconv.Atoi(ipStart[3])
	ipEnd0, _ := strconv.Atoi(ipEnd[0])
	ipEnd1, _ := strconv.Atoi(ipEnd[1])
	ipEnd2, _ := strconv.Atoi(ipEnd[2])
	ipEnd3, _ := strconv.Atoi(ipEnd[3])
	count = (ipEnd0 - ipStart0)*16777216 + (ipEnd1 - ipStart1)*65536+ (ipEnd2 - ipStart2)*256 + (ipEnd3 - ipStart3)+1
	return
}

func ReadConf()(config *conf.ConfigFile,err error){
	config,err = conf.ReadConfigFile("config.conf")
	if err!=nil{
		config,err = conf.ReadConfigFile("/etc/ccwings/canon/config.conf")
		if err !=nil{
			log.Error(err.Error())
		}
	}
	return
}