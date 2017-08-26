package Parser
import (
	"github.com/go-ini/ini"
	_"fmt"
	"fmt"
)
var zbxKey2App map[string]string =  make(map[string]string)
var app2itemNum map[string]int = make(map[string]int)


func getSections() (map[string]string,map[string]string){
	cfg, err := ini.InsensitiveLoad("./item.conf")

	if err != nil{
		panic("can not find items")
	}

	p2a := make(map[string]string)
	secs := cfg.Sections()
	itemKv := make(map[string]string)
	for _,s := range secs{
		keys := s.KeysHash()
		app2itemNum[s.Name()] = len(keys)
		for k,v := range keys{
			//fmt.Println(k.Value())
			itemKv[k] = v
			p2a[v] = k
			zbxKey2App[k] = s.Name()
		}
	}
	return itemKv,p2a
}

func QuerySection(query string) (string,error){
	fmt.Println("zbxKey2App",zbxKey2App)
	if zbxKey2App == nil{
		return "",fmt.Errorf("Empty Config Map")
	}else{
		elem, ok := zbxKey2App[query]
		if ok{
			return elem,nil
		}else{
			return zbxKey2App["REG_"+query],nil
		}

	}
}

func QuerySectionNumber(query string) int{
	return app2itemNum[query]
}

func ReadItems() (map[string]string,map[string]string) {
	return  getSections()
}

