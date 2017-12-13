package model

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/kusora/dlog"
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
	UserId            int64  `json:"user_id"`
	Feedback          string  `json:"feedback"`
	Status            int  `json:"status"`
	ServiceName       string `json:"service_name"`
	RelatedFeedbackId int64 `json:"related_feedback_id"`
	FeedbackType      int `json:"feedback_type"`
	CreatedAt         time.Time `json:"created_at"`
}

//func (fb *Feedback) MarshalJSON() ([]byte, error) {
//	m := map[string]interface{} {
//		"feedback_id": fb.FeedbackId,
//		"feedback": fb.Feedback,
//		"status": fb.Status,
//		"service_name": fb.ServiceName,
//		"feedback_type": fb.FeedbackType,
//		"created_at": fb.CreatedAt.Format("20060102150405"),
//	}
//
//	return json.Marshal(m)
//}


func (m *Model) InsertFeedback(fb *Feedback) error {
	return m.DoTransaction(func(o orm.Ormer) error {
		// 下面进行插入操作
		_, err := o.Insert(fb)
		return err
	})
}

func (m *Model) GetUserSendFeedbacks() ([]*Feedback, error) {
	var feedbacks []*Feedback
	cnt, err := m.m.QueryTable("feedback").Filter("feedback_type", 0).All(&feedbacks)
	if err != nil {
		return nil, err
	}
	dlog.Info("query %d records", cnt)
	return feedbacks, nil
}

func (m *Model) GetAllFeedbacks() ([]*Feedback, error) {
	var feedbacks []*Feedback
	_, err := m.m.QueryTable("feedback").All(&feedbacks)
	if err != nil {
		return nil, err
	}
	dlog.Info("query %d records", len(feedbacks))
	return feedbacks, nil
}