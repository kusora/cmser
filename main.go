package main

import (
	"github.com/kusora/dlog"
	"github.com/kusora/cmser/model"
	"os"
	"github.com/kusora/cmser/util"
	"github.com/kusora/cmser/cmd"
	"encoding/json"
	"io/ioutil"
	"fmt"
)
const (
	LEVEL_SAME = 0.60
)


func main() {
	m := model.NewModel()
	data, err := ioutil.ReadFile("conversation.json")
	if err != nil {
		dlog.Error("failed to read file, %+v", err)
		return
	}

	userConversation := make([]*cmd.UserConversation, 0)
	err = json.Unmarshal(data, &userConversation)
	if err != nil {
		dlog.Error("failed to unmarshal, %+v", err)
		return
	}

	qfs := make([]*model.Feedback, 0)
	// 这一遍先看看效果, 先不使用上下文关系，也不修改相似性函数
	for _, uc := range userConversation {
		for _, conversation := range uc.Conversations {
			for _, feedback := range conversation {
				if feedback.FeedbackType == model.FEEDBACK_TYPE_REPLY {
					qfs = append(qfs, feedback)
				}
			}
		}
	}

	groups := Groups(m, qfs)
	// 下面打印出相似



}


func ReadConversations(m *model.Model) {
	feedbacks, err := m.GetAllFeedbacks()
	if err != nil {
		dlog.Error("failed to get all feedbacks, err %+v", err)
		return
	}

	userConversations := cmd.GetConversations(feedbacks)
	data, err := json.Marshal(userConversations)
	if err != nil {
		dlog.Error("failed to marshal, err %+v", err)
		return
	}

	dlog.Info("%d, %d", len(userConversations), len(data))
	err = ioutil.WriteFile("conversations.json", data, 0666)
	if err != nil {
		dlog.Error("failed to write file %+v", err)
	}
}

// 设置阈值也是为了减少请求, 返回了分组
func Groups(m *model.Model, feedbacks []*model.Feedback) [][]int {
	// 这里按照分组来, 先处理1000条
	strs := make([]string, 0, 100000)
	for _, feedback := range feedbacks {
		strs = append(strs, feedback.Feedback)
	}

	groups := make([][]int, 0)
	flags := make(map[int]bool, 0)

	for id, str := range strs {
		if _, ok := flags[id]; ok {
			continue
		}
		rates, err := cmd.GetRelations(str, strs)
		if err != nil {
			dlog.Error("failed to get rates %+v", err)
			continue
		}
		group := make([]int, 0)
		for nid, rate := range rates {
			if rate > LEVEL_SAME {
				group = append(group, nid)
				flags[nid] = true
			}
		}
		if len(group) > 1 {
			groups = append(groups, group)
		}
	}
	return groups
}

//func readToDb(fname string) {
//	config.Init("./cmser.conf")
//	m := model.NewModel()
//	f, err := os.Open(fname)
//	if err != nil {
//		dlog.Error("%v", err)
//		os.Exit(1)
//	}
//
//	defer f.Close()
//
//	reader := bufio.NewReader(f)
//	line, err := reader.ReadString('\n')
//	dlog.Info("header: %s", line)
//	lineNum := 0
//	for {
//		line, err := reader.ReadString('\n')
//		parts := strings.Split(line, "\t")
//		//feedback_id user_id feedback status latitude longitude service_id service_name related_feedback_id feedback_type device_guid created_at platform app_version device_model user_agent
//		feedbackId, _ := strconv.ParseInt(parts[0], 10, 64)
//		userId, _ := strconv.ParseInt(parts[1], 10, 64)
//		feedback, _ := base64.StdEncoding.DecodeString(parts[2][7:len(parts[2])-1])
//		status, _ := strconv.Atoi(parts[3])
//		latitude, _ := strconv.ParseFloat(parts[4], 64)
//		longitude, _ := strconv.ParseFloat(parts[5], 64)
//		serviceId, _ := strconv.ParseInt(parts[6], 10, 64)
//		serviceName := parts[7]
//		relatedFeedbackId, _ := strconv.ParseInt(parts[8], 10, 64)
//		feedbackType, _ := strconv.Atoi(parts[9])
//		deviceGuid := parts[10]
//		createdAt, _ := time.ParseInLocation("2006-01-02 15:04:05", parts[11], time.Local)
//		platform := parts[12]
//		appVersion := parts[13]
//		deviceModel := parts[14]
//		userAgent := parts[15]
//
//
//		fb := &model.Feedback{
//			FeedbackId: feedbackId,
//			UserId: userId,
//			Feedback: string(feedback),
//			Status:  status,
//			//Latitude: latitude,
//			//Longitude:longitude,
//			//ServiceId       : serviceId,
//			ServiceName      : serviceName,
//			RelatedFeedbackId : relatedFeedbackId,
//			FeedbackType      : feedbackType,
//			//DeviceGuid        : deviceGuid,
//			CreatedAt         : createdAt,
//			//Platform          : platform,
//			//AppVersion        : appVersion,
//			//DeviceModel       :  deviceModel,
//			//UserAgent         : userAgent,
//		}
//
//		lineNum++
//		if err != nil {
//			dlog.Info("read %d lines", lineNum)
//			break
//		}
//
//		dlog.Info("message %d:%s %s.", lineNum, fb.Feedback, parts[11])
//		err = m.InsertFeedback(fb)
//		if err != nil {
//			dlog.Error("failed to insert feedback %v, %s", err, fb.Feedback)
//		}
//	}
//
//	dlog.Info("successfully insert %d feedback", lineNum)
//
//}
