package event

// MoreResp 过多的问题沉默或者回答正确
type MoreResp struct {
	Event      string   `json:"event"`
	BookID     string   `json:"bookId"`
	PageID     string   `json:"pageId"`
	UniqueID   string   `json:"uid"`
	QuestionID string   `json:"questionId"`
	Time       int64    `json:"time"` // 次数
	Count      int      `json:"count"`
	Questions  []string `json:"questions"`
	Type       string   `json:"continuousType"` // silence, right, wrong
	Mode       string   `json:"mode"`
}

// QuestionResp 问题的回答
type QuestionResp struct {
	Event        string   `json:"event"`
	BookID       string   `json:"bookId"`
	PageID       string   `json:"pageId"`
	UniqueID     string   `json:"uid"`
	Time         int64    `json:"time"`
	QuestionID   string   `json:"questionId"`
	ResponseTime float64  `json:"responseTime"`
	AnswerTimes  uint     `json:"answerTimes"`
	STTResult    string   `json:"sttResult"`
	Award        uint     `json:"award"`
	KeywordList  []string `json:"keywordList"`
	Mode         string   `json:"mode"`
}
