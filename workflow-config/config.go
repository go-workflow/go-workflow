package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mumushuiding/util"
)

// Configuration 数据库配置结构
type Configuration struct {
	Port           string
	ReadTimeout    string
	WriteTimeout   string
	DbLogMode      string
	DbType         string
	DbName         string
	DbHost         string
	DbPort         string
	DbUser         string
	DbPassword     string
	DbMaxIdleConns string
	DbMaxOpenConns string
}

// Config 数据库配置
var Config = &Configuration{}

func init() {
	LoadConfig()
}

// LoadConfig LoadConfig
func LoadConfig() {
	// 获取配置信息config
	Config.getConf()
	// 环境变量覆盖config
	err := Config.setFromEnv()
	if err != nil {
		panic(err)
	}
	// 打印配置信息
	config, _ := json.Marshal(&Config)
	log.Printf("configuration:%s", string(config))
}
func (c *Configuration) setFromEnv() error {
	// // 覆盖server配置
	// c.setServerFromEnv()
	// // 覆盖Db配置
	// c.setDbFromEnv()'
	// 获取对象Configuration的属性string流
	fieldStream, err := util.GetFieldChannelFromStruct(&Configuration{})
	if err != nil {
		return err
	}
	for fieldname := range fieldStream {
		if len(os.Getenv(fieldname)) > 0 {
			err = util.StructSetValByReflect(c, fieldname, os.Getenv(fieldname))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Configuration) getConf() *Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		log.Printf("cannot open file config.json：%v", err)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Printf("decode config.json failed:%v", err)
		panic(err)
	}
	return c
}
