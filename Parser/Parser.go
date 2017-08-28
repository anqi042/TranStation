package Parser

import "fmt"
import (
	//"encoding/json"
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
							//MyLogger.Info("get map:",zap.Any("hmap",m))
							//res := FormatOutput(m)
							//MyLogger.Info("formatted:",zap.String("result",res))
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
/*
check if the key is what we want;
if it equal to a string,then return its topic
if it match a reg,return its topic and match value
else return empty string

 */
func translateZbxKey(zbxkey string,regs []*regexp.Regexp, eqSet *hashset.Set,p2a map[string]string) (string,string){

	if eqSet.Contains(zbxkey){
		return p2a[zbxkey],""
	}
	for _,v := range regs{

		match := v.FindStringSubmatch(zbxkey)


		//MyLogger.Info("match info",zap.String(zbxkey,v.String()),zap.Any("match",len(match)))

		if len(match) == 2{
			return p2a[v.String()],match[1]
		}
	}

	return "",""

}

func parseBody(body []byte,inMap  chan *HostKV , regs []*regexp.Regexp, eqSet *hashset.Set, p2a map[string]string){
	const TAG_SPLIT_BODY="$"
	strs := strings.Split(string(body),TAG_SPLIT_BODY)
	data := make( []map[string]string,0)
	if len(strs) % 3 != 0{
		MyLogger.Info("BODY FUCKER ",zap.Any("this len",len(strs)))

		return
	}else {
		datai := 0
		//host key value
		for i := 0;i < len(strs) ;  {
			if strings.Contains(strs[i+1],"discovery") || strs[i] == ""||strs[i+1] == ""||strs[i+2] == ""{
				 //skip this group
			}else {
				data = append(data,make(map[string]string))
				data[datai]["host"] = strs[i]
				data[datai]["key"] = strs[i+1]
				data[datai]["value"] = strs[i+2]
				datai++

			}
			i += 3

		}
	}
	//debug
	//for _,_v := range regs{
	//	MyLogger.Info("reg",zap.String("reg",_v.String()))
	//}

	for _,v := range data{


		//filter here
		if str,match := translateZbxKey(v["key"],regs,eqSet,p2a);str != ""{
			//MyLogger.Info("Got a metric:",zap.String("topic",str))
			app, _ := QuerySection(str)
			if match != ""{
				tmp := &HostKV{v["host"], app +".TAG", match,app}
				MyLogger.Info("Got a tag:",zap.String("tag",match))
				inMap <- tmp
			}
			tmp := &HostKV{v["host"],str, v["value"],app}

			if len(inMap) > 900{
				//
				MyLogger.Warn("body to json channel is almost full")
				time.Sleep(time.Second)

			}else{
				MyLogger.Info("Add new hostKV",zap.Any("key",tmp.key),zap.Any("val",tmp.val),zap.Any("host",tmp.host))
				inMap <- tmp
			}


		}else{
			//MyLogger.Info("useless key, discard it")
		}

	}




}