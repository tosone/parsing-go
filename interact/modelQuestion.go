package interact

// questionSetting ..
type questionSetting struct {
	VeryBefore        questionSettingVeryBefore        `json:"veryBefore"`
	Before            questionSettingBefore            `json:"before"`
	AgainBefore       questionSettingAgainBefore       `json:"againBefore"`
	SaidNothing       questionSettingSaidNothing       `json:"saidNothing"`
	SaidWrong         questionSettingSaidWrong         `json:"saidWrong"`
	NoRecord          bool                             `json:"noRecord"`
	Sound             []string                         `json:"sound"`
	Hash              []string                         `json:"hash"`
	MaxRecording      int                              `json:"maxRecording"`
	Keyword           []string                         `json:"keyword"`
	Award             questionSettingAward             `json:"award"`
	RetryAsk          questionSettingRetryAsk          `json:"retryAsk"`
	NeverAgain        questionSettingNeverAgain        `json:"neverAgain"`
	NeverAgainNothing questionSettingNeverAgainNothing `json:"neverAgainNothing"`
	NeverAgainWrong   questionSettingNeverAgainWrong   `json:"neverAgainWrong"`
}

// questionSettingVeryBefore ..
type questionSettingVeryBefore struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingBefore ..
type questionSettingBefore struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingAgainBefore ..
type questionSettingAgainBefore struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingNeverAgainNothing ..
type questionSettingNeverAgainNothing struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingSaidNothing ..
type questionSettingSaidNothing struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingSaidWrong ..
type questionSettingSaidWrong struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingNeverAgainWrong ..
type questionSettingNeverAgainWrong struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingNeverAgain ..
type questionSettingNeverAgain struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingAward ..
type questionSettingAward struct {
	Num   int      `json:"num"`
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// questionSettingRetryAsk ..
type questionSettingRetryAsk struct {
	Num int `json:"num"`
}
