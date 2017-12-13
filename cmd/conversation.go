package cmd

import (
	"github.com/kusora/cmser/model"
	"sort"
	"github.com/kusora/dlog"
	"github.com/kusora/cmser/util"
)

type Reply struct {
	Ask string `json:"ask"`
	Answer string `json:"answer"`
}


type UserConversation struct {
	UserId int64 `json:"user_id"`
	Conversations [][]*model.Feedback `json:"conversations"`
}



/* 目标是找到有效对话
0. 生成对话
1. 对话分解，根据间隔时间以及对话最后的状态，但这里需要注意的是我们只标记了用户的问题的最终状态，所以会话不能忘记客服的最后一句话
2. 过滤掉没有解决问题的自动回复
3. 过滤掉用户的重复问题和无效问答
4. 给每个会话打分，
 */

func GetConversations(feedbacks []*model.Feedback) []*UserConversation {


	sort.Slice(feedbacks, func(i, j int) bool {
		if feedbacks[i].UserId < feedbacks[j].UserId {
			return true
		} else if feedbacks[i].UserId == feedbacks[j].UserId && feedbacks[i].FeedbackId < feedbacks[j].FeedbackId {
			return true
		}
		return false

	})

	// 下面过滤掉不合法的回答
	fbs := make([]*model.Feedback, 0)
	for id, feedback := range feedbacks {
		if !util.ContainChineseChar(feedback.Feedback) {
			continue
		}

		if id != 0 && feedbacks[id].UserId == feedbacks[id-1].UserId && feedbacks[id].Feedback == feedbacks[id-1].Feedback {
			continue
		}

		if len(feedback.Feedback) < 5 {
			continue
		}
		fbs = append(fbs, feedback)
	}

	// 下面开始进行处理
	result := make([]*UserConversation, 0)
	from := 0
	for id, feedback := range fbs {
		if id == 0 {
			continue
		}
		if feedback.UserId != fbs[id-1].UserId {
			conversation := BuildConversations(fbs[from:id])
			if conversation != nil {
				result = append(result, conversation)
			}
			from = id
		}
	}
	conversation := BuildConversations(fbs[from:])
	if conversation != nil {
		result = append(result, conversation)
	}

	return result
}

/*
	构建一个用户的会话记录
 */
func BuildConversations(feedbacks []*model.Feedback) *UserConversation {
	if len(feedbacks) == 0 {
		return nil
	}
	var uc UserConversation
	uc.UserId = feedbacks[0].UserId
	from := 0
	for id, feedback := range feedbacks {
		if feedback.UserId != uc.UserId {
			dlog.Error("wrong code, user_id should be same, %d, %d", uc.UserId, feedback.UserId)
			continue
		}
		if id == len(feedbacks) - 1 || (feedbacks[id+1].CreatedAt.Sub(feedback.CreatedAt).Hours() > 24) {
			if id + 1 - from <= 1 {
				continue
			}

			hasReply := false
			firstId := -1
			for id, feedback := range feedbacks[from:id+1] {
				if firstId < 0 && feedback.FeedbackType != 0 {
					continue
				}
				if feedback.FeedbackType == 0 && firstId < 0 {
					firstId = id
				}
				if feedback.FeedbackType == 1 {
					hasReply = true
					break
				}
			}
			if hasReply {
				uc.Conversations = append(uc.Conversations, feedbacks[from:id+1])
			}
			from = id + 1
		}
	}

	if len(uc.Conversations) == 0 {
		return nil
	}
	return &uc
}