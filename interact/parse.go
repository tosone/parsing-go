package interact

import (
	"encoding/json"
)

// parse 序列化传入的脚本
func (s *State) parse(interactScript []byte) (globalSetting globalSetting, questions map[string]map[string][]questionSetting, err error) {
	var data = map[string]json.RawMessage{}
	if err = json.Unmarshal(interactScript, &data); err != nil {
		return
	}

	var question = map[string][]questionSetting{}
	questions = map[string]map[string][]questionSetting{}
	for k, v := range data {
		if k == "globalSetting" {
			if err = json.Unmarshal(v, &globalSetting); err != nil {
				return
			}
		} else {
			if err = json.Unmarshal(v, &question); err != nil {
				return
			}
			questions[k] = question
		}
	}

	if globalSetting.MaxRecording == 0 {
		globalSetting.MaxRecording = 8
	}
	if globalSetting.RetryAsk.Num == 0 {
		globalSetting.RetryAsk.Num = 2
	}

	return
}
