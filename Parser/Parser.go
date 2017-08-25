package Parser

import "fmt"
import (
	"encoding/json"


	"strings"
	"github.com/emirpasic/gods/sets/hashset"
	"regexp"
)

type HostKV struct {
	host string
	key string
	val string
}



func GetItems(src map[string]string,m map[string]string)  {

}

func GenFilter(src map[string]string) ([]*regexp.Regexp,*hashset.Set) {
	var regs []*regexp.Regexp
	eqSet := hashset.New()

	for k,_ := range src{
		if strings.HasPrefix(k,"REG_"){

			r, _ := regexp.Compile(k)
			regs = append(regs,r)
		}else {
			eqSet.Add(k)
		}
	}
	return regs,eqSet
}

func RaiseParser(in chan []byte)  {
	/*
	const STATS_GOOD  = "TS_READY"
	const STATS_NOTGOOD  = "TS_NOTREADY"
	const STATS_KEY = "AREUREADY?"
*/
	itemKv := ReadItems()
	regs, eqSet := GenFilter(itemKv)

	host2items := make(map[string]map[string]string)


	mapC := make(chan *HostKV,1000)
	go func(_in chan []byte) {
		for{
			select {
			case byts := <- _in:
				//parse http body here
				parseBody(byts,mapC,regs,eqSet)
			}
		}


	}(in)


	go func() {
		for{
			select {
			case hkv := <- mapC:

				if host2items[hkv.host] == nil{
					//initialize map
					host2items[hkv.host] = make(map[string]string)
					host2items[hkv.host]["host"] = hkv.host
//					host2items[hkv.host][STATS_KEY] = STATS_GOOD
				}
				host2items[hkv.host][hkv.key] = hkv.val

				//GetItems(itemKv,host2items[hkv.host])
				fmt.Println(hkv.host,host2items[hkv.host][hkv.key])



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

func parseBody(body []byte,inMap  chan *HostKV , regs []*regexp.Regexp, eqSet *hashset.Set){
	var dat map[string]interface{}
	err := json.Unmarshal(body,&dat)
	if err != nil{
		fmt.Println(err)
	}else{
		data := dat["data"].([]interface{})
		for _,v := range data{
			theMap := v.(map[string]interface{})
			fmt.Println(theMap["host"], theMap["key"], theMap["value"])
			//filter here
			if eqSet.Contains(theMap["key"].(string)) || matchReg(theMap["key"].(string),regs){
				tmp := &HostKV{theMap["host"].(string),theMap["key"].(string), theMap["value"].(string)}
				inMap <- tmp
			}else{
				fmt.Println("drop it")
			}



		}
	}
}