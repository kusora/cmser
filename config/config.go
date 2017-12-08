package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/kusora/dlog"
)

type Config struct {
	Env           string
	Redis         Redis
	Logger        string
	DbConn        string
}


type Redis struct {
	RedisMasterName    string
	RedisSentinelAddrs []string
	RedisAddr          string // ip:port
	RedisPort          int
	RedisPwd           string
	RedisUser          string
}

var conf *Config

func Init(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		dlog.Error("failed to load config file %s, err %+v", path, err)
	}
	conf = &Config{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		dlog.Error("failed to unmarshal config file %s, err %+v", path, err)
	}
}

func Instance() *Config {
	if conf == nil {
		Init("./cmser.conf")
	}
	return conf
}

func GetCurrentDirectory() string {
	dlog.Info("heheheheheh")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dlog.Fatal("%+v", err)
	}

	dlog.Info("heheheheheh %s", dir)
	dlog.Info("heheheheheh")
	return strings.Replace(dir, "\\", "/", -1)
}
