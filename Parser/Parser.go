package Parser

import "fmt"
import (
	"encoding/json"
"github.com/mitchellh/copystructure"
)

type HostKV struct {
	host string
	key string
	val string
}

func RaiseParser(in chan []byte)  {
	host2items := make(map[string]map[string]string)

	mapReadyC := make(chan map[string]string)
	mapC := make(chan *HostKV,1000)
	go func(_in chan []byte) {
		for{
			select {
			case byts := <- _in:
				//parse http body here
				parseBody(byts,mapC)
			}
		}


	}(in)


	go func() {
		for{
			select {
			case hkv := <- mapC:

				if host2items[hkv.host] == nil{
					host2items[hkv.host] = make(map[string]string)
				}
				host2items[hkv.host][hkv.key] = hkv.val
				//fmt.Println(hkv.host,theMap[hkv.key])
				host2items[hkv.host]["host"] = hkv.host
				dup,err := copystructure.Copy(host2items[hkv.host])
				if err{
					fmt.Println(err)
				}else{
					mapReadyC <- dup.(map[string]string)
				}


			}
		}
	}()

	go func() {
		for{
			select {
			case kvs := <- mapReadyC:
				//determine whether we have what we want
				fmt.Println(kvs)
			}
		}

	}()


}


func parseBody(body []byte,inMap  chan *HostKV){
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
			tmp := &HostKV{theMap["host"].(string),theMap["key"].(string), theMap["value"].(string)}
			inMap <- tmp
		}
	}
}