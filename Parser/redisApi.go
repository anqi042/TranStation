package Parser

import (
	"github.com/go-redis/redis"
	"fmt"

	"strings"
)

var client *redis.Client =  redis.NewClient(&redis.Options{
Addr:     "localhost:6379",
Password: "", // no password set
DB:       0,  // use default DB
})


func PingRedis() {

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err != nil{
		panic(err)
	}
	// Output: PONG <nil>
}

func PushData(kv *HostKV){
	key := kv.host+":"+kv.app
	//delete any prefix
	k := strings.TrimPrefix(kv.key,"reg_")

	err := client.HSet(key,k,kv.val)
	err2 :=client.HSet(key,"host",kv.host)
	client.HSet(key,"topic",kv.app)
	if err != nil || err2 != nil{
		fmt.Println(err,err2)
	}
}

func AllKeys() ([]string,error){
	res := client.Keys("*")
	keys ,err := res.Result()
	return keys,err
}

func GetHMap(key string) map[string]string{
	fields := strings.Split(key,":")
	res := client.HLen(key)
	num,_ := res.Result()
	//number excludes  app and host
	num -= 2
	if num < int64(QuerySectionNumber(fields[1])){
		MyLogger.Info("not enough values")
		return nil
	}else{
		hmap := client.HGetAll(key)
		m := hmap.Val()
		return m
	}
}

