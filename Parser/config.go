package Parser
import (
	"github.com/go-ini/ini"
	_"fmt"
)

func getSections() map[string]string{
	cfg, err := ini.InsensitiveLoad("./item.conf")
	if err != nil{
		panic("can not find items")
	}

	secs := cfg.Sections()
	itemKv := make(map[string]string)
	for _,s := range secs{
		keys := s.KeysHash()
		for k,v := range keys{
			//fmt.Println(k.Value())
			itemKv[k] = v
		}
	}
	return itemKv
}
func ReadItems() map[string]string {
	return  getSections()
}

