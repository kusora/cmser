package cmd


type Reply struct {
	Message string `json:"message"`
	Weight  int `json:"weight"`  //这句话的用处
	From    string `json:"from"` //讲话者
	Ts      string `json:"ts"`   //20060102150405  讲话时间
}


type UserConversation struct {
	UserId int64 `json:"user_id"`
	Conversations [][]*Reply `json:"conversations"`
}

/* 目标是找到有效对话
0. 生成对话
1. 对话分解，根据间隔时间以及对话最后的状态，但这里需要注意的是我们只标记了用户的问题的最终状态，所以会话不能忘记客服的最后一句话
2. 过滤掉没有解决问题的自动回复
3. 过滤掉用户的重复问题和无效问答
4. 给每个会话打分，
 */
