package main

import (
	"github.com/kusora/dlog"
	"github.com/kusora/cmser/model"
	"os"
	"gopkg.in/bufio.v1"
	"strings"
	"github.com/kusora/cmser/config"
	"strconv"
	"encoding/base64"
	"time"
	"github.com/kusora/cmser/util"
	"github.com/kusora/cmser/cmd"
)
const (
	LEVEL_SAME = 0.75
)

func main() {
	m := model.NewModel()
	feedbacks, err := m.GetUserSendFeedbacks()
	if err != nil {
		os.Exit(1)
	}

	// 这里按照分组来, 先处理1000条
	strs := make([]string, 0, 100000)
	for _, feedback := range feedbacks {
		if !util.ContainChineseChar(feedback.Feedback) {
			dlog.Info("not contain chinese char %s", feedback.Feedback)
			continue
		}

		if len(feedback.Feedback) < 6 {
			continue
		}
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


	dlog.Info("has %d groups", len(groups))
	for id, group := range groups {
		dlog.Info("group: %d", id)
		for _, strId := range group {
			dlog.Info("%s", strs[strId])
		}
	}
}

func readToDb(fname string) {
	config.Init("./cmser.conf")
	m := model.NewModel()
	f, err := os.Open(fname)
	if err != nil {
		dlog.Error("%v", err)
		os.Exit(1)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	dlog.Info("header: %s", line)
	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		parts := strings.Split(line, "\t")
		//feedback_id user_id feedback status latitude longitude service_id service_name related_feedback_id feedback_type device_guid created_at platform app_version device_model user_agent
		feedbackId, _ := strconv.ParseInt(parts[0], 10, 64)
		userId, _ := strconv.ParseInt(parts[1], 10, 64)
		feedback, _ := base64.StdEncoding.DecodeString(parts[2][7:len(parts[2])-1])
		status, _ := strconv.Atoi(parts[3])
		latitude, _ := strconv.ParseFloat(parts[4], 64)
		longitude, _ := strconv.ParseFloat(parts[5], 64)
		serviceId, _ := strconv.ParseInt(parts[6], 10, 64)
		serviceName := parts[7]
		relatedFeedbackId, _ := strconv.ParseInt(parts[8], 10, 64)
		feedbackType, _ := strconv.Atoi(parts[9])
		deviceGuid := parts[10]
		createdAt, _ := time.ParseInLocation("2006-01-02 15:04:05", parts[11], time.Local)
		platform := parts[12]
		appVersion := parts[13]
		deviceModel := parts[14]
		userAgent := parts[15]


		fb := &model.Feedback{
			FeedbackId: feedbackId,
			UserId: userId,
			Feedback: string(feedback),
			Status:  status,
			Latitude: latitude,
			Longitude:longitude,
			ServiceId       : serviceId,
			ServiceName      : serviceName,
			RelatedFeedbackId : relatedFeedbackId,
			FeedbackType      : feedbackType,
			DeviceGuid        : deviceGuid,
			CreatedAt         : createdAt,
			Platform          : platform,
			AppVersion        : appVersion,
			DeviceModel       :  deviceModel,
			UserAgent         : userAgent,
		}

		lineNum++
		if err != nil {
			dlog.Info("read %d lines", lineNum)
			break
		}

		dlog.Info("message %d:%s %s.", lineNum, fb.Feedback, parts[11])
		err = m.InsertFeedback(fb)
		if err != nil {
			dlog.Error("failed to insert feedback %v, %s", err, fb.Feedback)
		}
	}

	dlog.Info("successfully insert %d feedback", lineNum)

}
