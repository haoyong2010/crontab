package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//程序配置
type Config struct {
	EtcdEndpoints         []string `json:"etcd_endpoints"`
	EtcdDialTimeout       int      `json:"etcd_dial_timeout"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
	JobLogBatchSize       int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout   int      `json:"jobLogCommitTimeout"`
}

var (
	//定义一个单例
	G_config *Config
)

func InitConfig(filename string) (err error) {
	fmt.Println(filename)
	var (
		content []byte
		conf    Config
	)
	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	//反序列化JSON
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	//赋值单例
	G_config = &conf
	fmt.Println(conf)
	return
}
