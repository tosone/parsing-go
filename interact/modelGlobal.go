package interact

// Pages ..
type Pages map[string]interface{}

// globalSetting ..
type globalSetting struct {
	QuestionsPkg           globalSettingQuestionsPkg           `json:"questionsPkg"`
	Version                string                              `json:"version"`
	ContinuousBackinWords  globalSettingContinuousBackinWords  `json:"continuousBackinWords"`
	ContinuousBackSilence  globalSettingContinuousBackSlience  `json:"continuousBackSlience"`
	ContinuousBackoutWords globalSettingContinuousBackoutWords `json:"continuousBackoutWords"`
	SaidNothing            globalSettingSaidNothing            `json:"saidNothing"`
	SaidWrong              globalSettingSaidWrong              `json:"saidWrong"`
	ReadOver               globalSettingReadOver               `json:"readOver"`
	RetryAsk               globalSettingRetryAsk               `json:"retryAsk"`
	MaxRecording           int                                 `json:"maxRecording"`
	PerPageQuestion        globalSettingPerPageQuestion        `json:"perPageQuestion"`
	NeverAgainNothing      globalSettingNeverAgainNothing      `json:"neverAgainNothing"`
	NeverAgainWrong        globalSettingNeverAgainWrong        `json:"neverAgainWrong"`
	NeverAgain             globalSettingNeverAgain             `json:"neverAgain"`
	BreakBy                string                              `json:"breakBy"`
}

// globalSettingNeverAgain ..
type globalSettingNeverAgain struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingNeverAgainNothing ..
type globalSettingNeverAgainNothing struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingNeverAgainWrong ..
type globalSettingNeverAgainWrong struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingQuestionsPkg ..
type globalSettingQuestionsPkg struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

// globalSettingContinuousBackoutWords ..
type globalSettingContinuousBackoutWords struct {
	Num   int      `json:"num"`
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingContinuousBackinWords ..
type globalSettingContinuousBackinWords struct {
	Num        int      `json:"num"`
	RetryReset bool     `json:"retryReset"`
	Sound      []string `json:"sound"`
	Hash       []string `json:"hash"`
}

// globalSettingContinuousBackSlience ..
type globalSettingContinuousBackSlience struct {
	Num   int      `json:"num"`
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingSaidNothing ..
type globalSettingSaidNothing struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingSaidWrong ..
type globalSettingSaidWrong struct {
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingReadOver ..
type globalSettingReadOver struct {
	Num   int      `json:"num"`
	Sound []string `json:"sound"`
	Hash  []string `json:"hash"`
}

// globalSettingRetryAsk ..
type globalSettingRetryAsk struct {
	Time int  `json:"time"` // 在多少秒之前没有说话认为是需要重新问问题，默认 4 秒
	Num  uint `json:"num"`  // 问题需要问多少遍，默认两遍
}

// globalSettingPerPageQuestion ..
type globalSettingPerPageQuestion struct {
	ProbabilityMode string `json:"probabilityMode"`
	Num             int    `json:"num"`
}
