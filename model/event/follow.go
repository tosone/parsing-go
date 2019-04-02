package event

// FollowRecordTakeOff 跟读过程中的事件内容体
type FollowRecordTakeOff struct {
	Event           string `json:"event"`
	InsertID        string `json:"insertId"`
	BookID          string `json:"bookId"`
	UniqueID        string `json:"uid"`
	Task            uint   `json:"task"`
	FinishTask      uint   `json:"finishTask"`
	ExtraTask       uint   `json:"extraTask"`
	FinishExtraTask uint   `json:"finishExtraTask"`
	Flower          uint   `json:"flower"`
	ValidFollow     uint   `json:"validFollow"`
}

// FollowRecordInfo 跟读录音上传
type FollowRecordInfo struct {
	Event    string                 `json:"event"`    // 要上传的时间内容
	InsertID string                 `json:"insertId"` // ..
	BookID   string                 `json:"bookId"`   // 书的 ID
	UniqueID string                 `json:"uid"`      // 书的版本的唯一指定
	PageID   string                 `json:"pageId"`   // 页的 ID
	FollowID string                 `json:"followId"` // 跟读 ID
	Start    int64                  `json:"start"`    // 时间戳，毫秒为单位
	End      int64                  `json:"end"`      // 时间戳，毫秒为单位
	Score    uint                   `json:"score"`    // 0 到 132 之间的参数
	Detail   FollowRecordInfoDetail `json:"detail"`   // 0 到 100 之间的数字
	Times    uint                   `json:"times"`    // 当前已经跟读的次数
	Record   string                 `json:"record"`   // 要上传的音频内容
}

// FollowRecordInfoDetail 详细信息
type FollowRecordInfoDetail struct {
	Correct uint `json:"correct"` // 准确率
	Speed   uint `json:"speed"`   // 速度
	Tone    uint `json:"tone"`    // 语调
}
