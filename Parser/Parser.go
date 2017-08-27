package Parser

import "fmt"
import (
	"encoding/json"
	"strings"
	"github.com/emirpasic/gods/sets/hashset"
	"regexp"
   "go.uber.org/zap"
	"time"

)

type HostKV struct {
	host string
	key string
	val string
	app string
}

const REDIS_GET_INTERVAL  = 500
const CHANNEL_HTTP_IN_LEN = 1000
const CHANNEL_TO_REDIS_LEN = 1000
const BATCH_SEND_SIZ = 2

func GetItems(src map[string]string,m map[string]string)  {

}

func GenFilter(src map[string]string) ([]*regexp.Regexp,*hashset.Set) {
	var regs []*regexp.Regexp
	eqSet := hashset.New()

	for k,v := range src{
		if strings.HasPrefix(k,"reg_"){

			r, err := regexp.Compile(v)
			if err != nil || r == nil{
				MyLogger.Error("can not compile regexp")
			}
			regs = append(regs,r)
		}else {
			eqSet.Add(v)
		}
	}
	fmt.Println("regs",regs,"strs",eqSet)
	return regs,eqSet
}

func RaiseParser(in chan []byte)  {

	itemKv,p2a := ReadItems()
	regs, eqSet := GenFilter(itemKv)
	PingRedis()



	mapC := make(chan *HostKV,CHANNEL_TO_REDIS_LEN)
	go func(_in chan []byte) {
		for{
			select {
			case byts := <- _in:
				//parse http body here
				parseBody(byts,mapC,regs,eqSet,p2a)
			}
		}


	}(in)


	go func() {
		for{
			select {
			case hkv := <- mapC:

			//	fmt.Println("host-k-v-App",hkv.host,hkv.key,hkv.val,hkv.app)
				PushData(hkv)
			}
		}
	}()


	//output
	go func() {
		ticker := time.NewTicker(time.Millisecond * REDIS_GET_INTERVAL)
		c := ticker.C
		for {
			select {
			case <- c:
				ks,e := AllKeys()
				if e == nil{
					for _,k := range ks{
						m := GetHMap(k)
						if m != nil{
							MyLogger.Info("get map:",zap.Any("hmap",m))

						}

					}
				}
			}
		}


	}()

}

func matchReg(src string,regs []*regexp.Regexp) bool {
	for _,v := range regs{
		if v.Match([]byte(src)){
			return true
		}
	}
	return false
}

func translateZbxKey(zbxkey string,regs []*regexp.Regexp, eqSet *hashset.Set,p2a map[string]string) (string,bool){

	if eqSet.Contains(zbxkey){
		return p2a[zbxkey],false
	}
	for _,v := range regs{
		if v.Match([]byte(zbxkey)){
			return p2a[v.String()],true
		}
	}

	return "",false

}

func parseBody(body []byte,inMap  chan *HostKV , regs []*regexp.Regexp, eqSet *hashset.Set, p2a map[string]string){
	var dat map[string]interface{}
	rNum,_ := regexp.Compile("[0-9]+")
	err := json.Unmarshal(body,&dat)
	if err != nil{
		fmt.Println(err)
	}else{

		data := dat["data"].([]interface{})
		for _,v := range data{
			theMap := v.(map[string]interface{})
			fmt.Println(theMap["host"], theMap["key"], theMap["value"])
			//filter here
			if str,needNumber := translateZbxKey(theMap["key"].(string),regs,eqSet,p2a);str != ""{
				fmt.Println(str)
				app, _ := QuerySection(str)
				if needNumber{
					res := rNum.Find([]byte(theMap["key"].(string)))
					fmt.Println("res Number=",string(res))

					tmp := &HostKV{theMap["host"].(string), app +".TAG", string(res),app}
					inMap <- tmp
				}
				tmp := &HostKV{theMap["host"].(string),str, theMap["value"].(string),app}

				if len(inMap) > 900{
					//
					MyLogger.Warn("body to json channel is almost full")
					time.Sleep(time.Second)

				}else{
					inMap <- tmp
				}


			}else{
				MyLogger.Info("useless key, discard it")
			}
		}
	}
}