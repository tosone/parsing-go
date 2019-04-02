package follow

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"smartconn.cc/tosone/parsing-go/logstash"
	"smartconn.cc/tosone/parsing-go/model"
	"smartconn.cc/tosone/parsing-go/tables"
	"smartconn.cc/tosone/parsing-go/util"
)

// randomSelectHash 传入的 hash 随机选择某个音频
func randomSelectHash(tips []model.Tip, hash []string) (string, string) {
	if len(hash) == 0 {
		return "", ""
	}
	h := util.RandSelect(hash)
	for _, audio := range tips {
		if h == audio.Hash {
			return audio.Path, h
		}
	}
	return "", ""
}

// randomBaseScript 开始这个模式之前选择一次 base 脚本
func (s State) randomBaseScript() (baseScript GlobalSetting, err error) {
	//var baseFollowScript []tables.FollowBaseScript
	//
	//if err = store.Major.Find(&baseFollowScript).Error; err != nil {
	//	return
	//}
	var baseFollowScript = []tables.FollowBaseScript{{Version: "", Manifest: []byte(s.baseScriptRaw), ScriptID: ""}}

	if len(baseFollowScript) == 0 {
		err = fmt.Errorf("cannot find any base script")
		return
	} else if len(baseFollowScript) == 1 {
		if err = json.Unmarshal(baseFollowScript[0].Manifest, &baseScript); err != nil {
			return
		}
		return
	} else {
		for {
			rand.Seed(time.Now().UTC().UnixNano())
			targetBaseScriptID := rand.Intn(len(baseFollowScript))
			if targetBaseScriptID == baseScriptID {
				continue
			} else {
				logstash.Info(fmt.Sprintf("Base script select id: %s", baseFollowScript[targetBaseScriptID].Version))
				baseScriptID = targetBaseScriptID
				if err = json.Unmarshal(baseFollowScript[targetBaseScriptID].Manifest, &baseScript); err != nil {
					return
				}
				break
			}
		}
	}

	return
}
