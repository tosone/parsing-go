package follow

import "encoding/json"

// Pages ..
type Pages map[string]json.RawMessage

// GlobalSetting 全局配置
type GlobalSetting struct {
	BaseVersion         string
	LocalVersion        string
	BasePkg             GlobalSettingPkg               `json:"basePkg"`
	FollowsPkg          GlobalSettingPkg               `json:"followsPkg"`
	Version             string                         `json:"version"`
	MaxFollowInvalid    *GlobalSettingMaxFollowInvalid `json:"maxFollowInvalid"`
	Goal                *Goal                          `json:"goal"`
	Progress            *[]GlobalSoundWithHash         `json:"progress"`
	ProgressComplete    *[]GlobalSoundWithHash         `json:"progressComplete"`
	Extra               *Goal                          `json:"extra"`
	ExtraProgress       *[]GlobalSoundWithHash         `json:"extraProgress"`
	Retry               *GlobalSettingRetry            `json:"retry"`
	NextPage            *GlobalSoundWithHash           `json:"nextPage"`
	Evaluation          *[]Evaluation                  `json:"evaluation"`
	IntroTask           *[]GlobalSoundWithHash         `json:"introTask"`
	ContentBefore       *GlobalSoundWithHash           `json:"contentBefore"`
	ContentAfterSuccess *GlobalSoundWithHash           `json:"contentAfterSucc"`
	ContentAfterFail    *[]GlobalSoundWithHash         `json:"contentAfterFail"`
	TakeOff             *GlobalSettingContentAfter     `json:"takeOff"`
	BeforeFollow        *GlobalSoundWithHash           `json:"beforeFollow"`
	ExcludeContent      []string                       `json:"excludeContent"`
	LastContentPage     string                         `json:"lastContentPage"`
	MaxRecording        *int                           `json:"maxRecording"`
	BeforeRecord        *GlobalSoundWithHash           `json:"beforeRecord"`
}

// GlobalSettingPkg ..
type GlobalSettingPkg struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

// GlobalSettingRetry ..
type GlobalSettingRetry struct {
	MaxRetry  uint                `json:"maxRetry"`
	BadScore  GlobalSoundWithHash `json:"badScore"`
	Condition string              `json:"condition"`
	Score     uint                `json:"score"`
}

// GlobalSettingMaxFollowInvalid ..
type GlobalSettingMaxFollowInvalid struct {
	Num   int      `json:"num"`
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// GlobalSoundWithHash ..
type GlobalSoundWithHash struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// GlobalSettingContentAfter ..
type GlobalSettingContentAfter struct {
	Fail    GlobalSoundWithHash `json:"fail"`
	Success GlobalSoundWithHash `json:"success"`
}

// Evaluation ..
type Evaluation struct {
	MaxScore uint     `json:"maxScore"`
	MinScore uint     `json:"minScore"`
	Valid    bool     `json:"valid"`
	Sound    []string `json:"sound"`
	Hash     []string `json:"hash"`
}

// Goal ..
type Goal struct {
	SentencePerFlower string `json:"sentencePerFlower"`
	FlowerPerTask     string `json:"flowerPerTask"`
	Task              int    `json:"task"`
}

// Progress ..
type Progress struct {
	Progress int      `json:"progress"`
	Sound    []string `json:"sound"`
	Hash     []string `json:"hash"`
}

// WholePages ..
type WholePages map[string]FollowItem

// FollowItem ..
type FollowItem struct {
	NoRecord   bool
	Order      []string
	InnerOrder map[string][]string
	Follows    map[string]FollowActon
}

// FollowSetting ..
type FollowSetting map[string]json.RawMessage

// FollowActon ..
type FollowActon struct {
	StandardOgg struct {
		Sound []string `json:"sound"`
		Hash  []string `json:"hash"`
	} `json:"standardOgg"`
	StandardWav struct {
		Sound []string `json:"sound"`
		Hash  []string `json:"hash"`
	} `json:"standardWav"`
}
