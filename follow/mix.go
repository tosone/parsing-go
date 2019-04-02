package follow

import (
	"fmt"

	"smartconn.cc/tosone/parsing-go/logstash"
)

var defaultMaxFollowInvalid = 3
var defaultMaxRecording = 8

// mixGlobal 将 base 脚本中的配置和当前这本书的 global setting 混合到一块
func mixGlobal(baseScript, localScript GlobalSetting) (mixed GlobalSetting, err error) {
	mixed.BasePkg = baseScript.BasePkg
	mixed.FollowsPkg = localScript.FollowsPkg

	mixed.BaseVersion = baseScript.Version
	mixed.LocalVersion = localScript.Version

	mixed.MaxFollowInvalid = &GlobalSettingMaxFollowInvalid{}

	if localScript.BeforeRecord != nil {
		mixed.BeforeRecord = localScript.BeforeRecord
	} else if baseScript.BeforeRecord != nil {
		mixed.BeforeRecord = baseScript.BeforeRecord
	}

	if localScript.MaxFollowInvalid != nil && len(localScript.MaxFollowInvalid.Hash) != 0 {
		mixed.MaxFollowInvalid = localScript.MaxFollowInvalid
	} else if baseScript.MaxFollowInvalid != nil && len(baseScript.MaxFollowInvalid.Hash) != 0 {
		mixed.MaxFollowInvalid = baseScript.MaxFollowInvalid
	}
	if localScript.MaxFollowInvalid != nil && localScript.MaxFollowInvalid.Num != 0 {
		mixed.MaxFollowInvalid.Num = localScript.MaxFollowInvalid.Num
	} else if baseScript.MaxFollowInvalid != nil && baseScript.MaxFollowInvalid.Num != 0 {
		mixed.MaxFollowInvalid.Num = baseScript.MaxFollowInvalid.Num
	} else {
		mixed.MaxFollowInvalid.Num = defaultMaxFollowInvalid
	}
	if localScript.Goal != nil {
		mixed.Goal = localScript.Goal
	} else if baseScript.Goal != nil {
		mixed.Goal = baseScript.Goal
	}
	if localScript.Progress != nil {
		mixed.Progress = localScript.Progress
	} else if baseScript.Progress != nil {
		mixed.Progress = baseScript.Progress
	}
	if localScript.ExtraProgress != nil {
		mixed.ExtraProgress = localScript.ExtraProgress
	} else if baseScript.ExtraProgress != nil {
		mixed.ExtraProgress = baseScript.ExtraProgress
	}
	if localScript.ProgressComplete != nil {
		mixed.ProgressComplete = localScript.ProgressComplete
	} else if baseScript.ProgressComplete != nil {
		mixed.ProgressComplete = baseScript.ProgressComplete
	}
	if localScript.Extra != nil {
		mixed.Extra = localScript.Extra
	} else if baseScript.Extra != nil {
		mixed.Extra = baseScript.Extra
	}
	if localScript.Retry != nil {
		mixed.Retry = localScript.Retry
	} else if baseScript.Retry != nil {
		mixed.Retry = baseScript.Retry
	} else {
		mixed.Retry.Score = 20
		mixed.Retry.Condition = "silence"
		mixed.Retry.MaxRetry = 2
	}
	if localScript.NextPage != nil {
		mixed.NextPage = localScript.NextPage
	} else if baseScript.NextPage != nil {
		mixed.NextPage = baseScript.NextPage
	}
	if localScript.IntroTask != nil {
		mixed.IntroTask = localScript.IntroTask
	} else if baseScript.IntroTask != nil {
		mixed.IntroTask = baseScript.IntroTask
	}
	if localScript.ContentBefore != nil {
		mixed.ContentBefore = localScript.ContentBefore
	} else if baseScript.ContentBefore != nil {
		mixed.ContentBefore = baseScript.ContentBefore
	}
	if localScript.ContentAfterSuccess != nil {
		mixed.ContentAfterSuccess = localScript.ContentAfterSuccess
	} else if baseScript.ContentAfterSuccess != nil {
		mixed.ContentAfterSuccess = baseScript.ContentAfterSuccess
	}
	if localScript.ContentAfterFail != nil {
		mixed.ContentAfterFail = localScript.ContentAfterFail
	} else if baseScript.ContentAfterFail != nil {
		mixed.ContentAfterFail = baseScript.ContentAfterFail
	}
	if localScript.BeforeFollow != nil {
		mixed.BeforeFollow = localScript.BeforeFollow
	} else if baseScript.BeforeFollow != nil {
		mixed.BeforeFollow = baseScript.BeforeFollow
	}
	if localScript.MaxRecording != nil {
		mixed.MaxRecording = localScript.MaxRecording
	} else if baseScript.MaxRecording != nil {
		mixed.MaxRecording = baseScript.MaxRecording
	} else {
		mixed.MaxRecording = &defaultMaxRecording
	}
	if localScript.Evaluation != nil {
		mixed.Evaluation = localScript.Evaluation
	} else if baseScript.Evaluation != nil {
		mixed.Evaluation = baseScript.Evaluation
	}
	if localScript.Progress != nil {
		mixed.Progress = localScript.Progress
	} else if baseScript.Progress != nil {
		mixed.Progress = baseScript.Progress
	}
	if localScript.ExtraProgress != nil {
		mixed.ExtraProgress = localScript.ExtraProgress
	} else if baseScript.ExtraProgress != nil {
		mixed.ExtraProgress = baseScript.ExtraProgress
	}
	if localScript.TakeOff != nil {
		mixed.TakeOff = localScript.TakeOff
	} else if baseScript.TakeOff != nil {
		mixed.TakeOff = baseScript.TakeOff
	}
	mixed.ExcludeContent = localScript.ExcludeContent
	mixed.LastContentPage = localScript.LastContentPage
	logstash.Info(fmt.Sprintf("baseScript: %+v\nlocalScript:%+v\nmixed: %+v",
		baseScript, localScript, mixed))
	return
}
