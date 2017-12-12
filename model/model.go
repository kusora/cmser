package model

import (
	"github.com/astaxie/beego/orm"
	"errors"
	"gopkg.in/redis.v3"
	"time"
	"sync"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	//"github.com/jiongzhao/kuaicha/cache"
	"github.com/kusora/dlog"
	//"reflect"
	//"github.com/jiongzhao/order-trus/util"
	"github.com/kusora/cmser/config"
)

type Model struct {
	m           orm.Ormer
	redisClient *redis.Client
	//cache       *cache.SimpleCache
}

var tables []interface{}
func init() {
	tables = []interface{}{
		new(Feedback),
	}
}

func NewModel() *Model {
	once.Do(func() {
		orm.RegisterDriver("mysql", orm.DRMySQL)
		// orm.DefaultTimeLoc = time.UTC
		orm.RegisterDataBase("default", "mysql", config.Instance().DbConn, 30/*maxIdle*/, 30/* maxConn*/)
		orm.RegisterModel(tables...)
		orm.DefaultRowsLimit = 1000000
	})

	var redisClient *redis.Client = nil
	redisPwd := ""
	if config.Instance().Redis.RedisUser != "" {
		redisPwd = fmt.Sprintf("%s:%s", config.Instance().Redis.RedisUser, config.Instance().Redis.RedisPwd)
	}

	//如果单机版的redis配置存在, 那么就走单机版路线
	if config.Instance().Redis.RedisAddr != "" {
		opt := &redis.Options{
			Addr: fmt.Sprintf("%s:%d", config.Instance().Redis.RedisAddr, config.Instance().Redis.RedisPort),
			Password: redisPwd,
			PoolSize:      10,
			DialTimeout:   time.Second,
			ReadTimeout:   time.Second,
			WriteTimeout:  time.Second,
			IdleTimeout:   time.Second * 10,
		}

		redisClient = redis.NewClient(opt)
	} else if config.Instance().Redis.RedisMasterName != "" {
		failoverOpt := &redis.FailoverOptions{
			MasterName:    config.Instance().Redis.RedisMasterName,
			SentinelAddrs: config.Instance().Redis.RedisSentinelAddrs,
			PoolSize:      10,
			DialTimeout:   time.Second,
			ReadTimeout:   time.Second,
			WriteTimeout:  time.Second,
			IdleTimeout:   time.Second * 10,
		}
		redisClient = redis.NewFailoverClient(failoverOpt)
	}

	if redisClient == nil {
		panic("redis setting error")
	}

	if redisClient.Ping().Val() != "PONG" {
		panic("redis ping result is not PONG")
	}

	m := &Model{
		m : orm.NewOrm(),
		redisClient:   redisClient,
		//cache: cache.NewSimpleCache(),
	}
	return m
}

var once sync.Once

func (m *Model) TruncateTables() error {
	//if config.Instance().Env == "staging" {
	//	m.redisClient.FlushAll()
	//	for _, table := range tables {
	//		t := reflect.TypeOf(table)
	//		tableName := util.Camel2underscore(t.Elem().Name())
	//		result := m.m.Raw("truncate table `" + tableName + "`")
	//		_, err := result.Exec()
	//		if err != nil {
	//			fmt.Println(err)
	//			return err
	//		}
	//	}
	//	return nil
	//}
	return errors.New("not allow to truncate table other than staging env")
}

func (m *Model) Shutdown() {
	// TODO close db pool
	err := m.redisClient.Close()
	if err != nil {
		dlog.Error("failed to close slave orm", err)
	}
}

func (m *Model) DoTransaction(f func(O orm.Ormer) error) error {
	// using transactions will affect Ormer object and all its queries.
	// So don’t use a global Ormer object if you need to switch databases or use transactions.
	o := orm.NewOrm()
	o.Begin()
	err := f(o)
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	return err
}

/*
Exists

检查redis中key是否存在，不存在的话就进行设置，返回false， 存在的话返回true
 */
func (m *Model) CheckAndSet(key string, duration time.Duration) bool {
	result := m.redisClient.Exists(key)
	if result.Err() == redis.Nil || !result.Val() {
		m.redisClient.Set(key, "1", duration)
		return false
	}
	return true
}

func (m *Model) Redis() *redis.Client {
	return m.redisClient
}


