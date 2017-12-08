package model

import (
	"time"
	"github.com/astaxie/beego/orm"
)

const (
	FEEDBACK_TYPE_USER       = 0
	FEEDBACK_TYPE_REPLY      = 1
	FEEDBACK_STATUS_WAITING  = 0
	FEEDBACK_STATUS_PROCESS  = 1
	FEEDBACK_STATUS_COMPLETE = 2
)

type Feedback struct {
	FeedbackId        int64 `json:"feedback_id" orm:"pk"`
	UserId            int64 //0: 表示该用户匿名
	Feedback          string
	Status            int
	Latitude          float64
	Longitude         float64
	ServiceId         int64 //不再使用ServiceId，采用ServiceName来表示客服身份
	ServiceName       string
	RelatedFeedbackId int64
	FeedbackType      int // 0: 用户上行 1：客服回复
	DeviceGuid        string
	CreatedAt         time.Time
	Platform          string
	AppVersion        string
	DeviceModel       string
	UserAgent         string
}


func (m *Model) InsertFeedback(fb *Feedback) error {
	return m.DoTransaction(func(o orm.Ormer) error {
		// 下面进行插入操作
		_, err := o.Insert(fb)
		return err
	})
}