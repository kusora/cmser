package main

import (
	"github.com/kusora/dlog"
	"github.com/kusora/cmser/model"
	"github.com/kusora/cmser/cmd"
	"encoding/json"
	"io/ioutil"
	"time"
	"sort"
)

const (
	LEVEL_SAME = 0.70
)


func main() {
	data, err := ioutil.ReadFile("groups_70.json")
	if err != nil {
		dlog.Error("failed to read file %+v", err)
		return
	}

	result := make([]interface{}, 0)

	err = json.Unmarshal(data, &result)
	if err != nil {
		dlog.Error("failed to unmarshao %+v", err)
		return
	}

	newData, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		dlog.Error("failed to marshal %+v", err)
		return
	}

	ioutil.WriteFile("groups_indent.json", newData, 0666)
}

func SaveAllGorups() {
	m := model.NewModel()
	//SaveConversations(m)
	data, err := ioutil.ReadFile("./conversations.json")
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
	mfs := make(map[int64]*model.Feedback, 0)
	u2cs := make(map[int64]*cmd.UserConversation, 0)
	// 这一遍先看看效果, 先不使用上下文关系，也不修改相似性函数
	for _, uc := range userConversation {
		u2cs[uc.UserId] = uc
		for _, conversation := range uc.Conversations {
			for _, feedback := range conversation {
				mfs[feedback.FeedbackId] = feedback
				if feedback.FeedbackType == model.FEEDBACK_TYPE_REPLY {
					qfs = append(qfs, feedback)
				}
			}
		}
	}

	groups := Groups(m, qfs)
	result := make([][][]string, 0)
	for _, group := range groups {
		groupInfo := make([][]string, 0)
		// 先收集问题
		asks := make([]string, 0)
		answers := make([]string, 0)
		for _, q := range group {
			feedback := mfs[q]
			answers = append(answers, feedback.Feedback)
			if feedback.RelatedFeedbackId > 0 {
				if fb, ok := mfs[feedback.RelatedFeedbackId]; ok {
					asks = append(asks, fb.Feedback)
					continue
				}
			}
			uc := u2cs[feedback.UserId]
			for _, conversation := range uc.Conversations {
				for id, fb := range conversation {
					if fb.FeedbackId == feedback.FeedbackId && id > 0 && conversation[id-1].FeedbackType == model.FEEDBACK_TYPE_USER {
						asks = append(asks, conversation[id-1].Feedback)
					}
				}
			}
		}
		asks = Dedup(asks)
		answers = Dedup(answers)
		groupInfo = append(groupInfo, asks, answers)
		result = append(result, groupInfo)
	}

	groupData, _ := json.Marshal(result)
	ioutil.WriteFile("groups.json", groupData, 0666)
}


func Dedup(strs []string) []string {
	if len(strs) == 0 {
		return strs
	}

	sort.Strings(strs)
	newStrs := []string{strs[0]}
	for _, str := range strs {
		if str != newStrs[len(newStrs)-1] {
			newStrs = append(newStrs, str)
		}
	}
	return newStrs
}

func SaveConversations(m *model.Model) {
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
func Groups(m *model.Model, feedbacks []*model.Feedback) [][]int64 {
	dlog.Info("has %d feedbacks", len(feedbacks))
	// 这里按照分组来, 先处理1000条
	strs := make([]string, 0, 100000)
	for _, feedback := range feedbacks {
		strs = append(strs, feedback.Feedback)
	}

	groups := make([][]int64, 0)
	flags := make(map[int]bool, 0)
	leftStrs := strs
	leftIds := make([]int, 0)
	for id, _ := range strs {
		leftIds = append(leftIds, id)
	}

	for id, str := range strs {
		now := time.Now()
		if id % 30 == 0 {
			leftStrs = make([]string, 0)
			leftIds = make([]int, 0)
			for newid, newstr := range strs {
				if _, ok := flags[newid]; ok {
					continue
				}
				leftStrs = append(leftStrs, newstr)
				leftIds = append(leftIds, newid)
			}
		}


		if _, ok := flags[id]; ok {
			continue
		}


		rates, err := cmd.GetRelations(str, leftStrs)
		if err != nil {
			dlog.Error("failed to get rates %+v, %d", err, id)
			continue
		}
		group := make([]int64, 0)
		for nid, rate := range rates {
			if rate > LEVEL_SAME {
				group = append(group, feedbacks[leftIds[nid]].FeedbackId)
				flags[leftIds[nid]] = true
			}
		}
		if len(group) > 1 {
			groups = append(groups, group)
		}
		dlog.Info("finish query %d,  %.2f seconds, groups %d, passed fbs %d", id, time.Now().Sub(now).Seconds(), len(groups), len(flags))
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
