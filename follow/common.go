package follow

import (
	"fmt"

	"github.com/Unknwon/com"
	"smartconn.cc/tosone/parsing-go/errs"
	"smartconn.cc/tosone/parsing-go/logstash"
	"smartconn.cc/tosone/parsing-go/model"
	"smartconn.cc/tosone/parsing-go/util"
)

func (s *State) playPureWithCtx(sound string) (err error) {
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	if err = s.audio.Play(sound); err != nil {
		return
	}
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	return
}

const recordCut = 32 * 200

func (s *State) playPCMWithCtx(sound string) (err error) {
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	if !com.IsFile(sound) {
		return fmt.Errorf("No such a file: %s", sound)
	}
	if err = s.audio.Play(sound); err != nil {
		return
	}
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	return
}

func (s *State) playWithCtx(sound, hash []string, action ...FollowActon) (err error) {
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	audioFile, audioHash := randomSelectHash(s.tips, hash)
	if !com.IsFile(audioFile) {
		err = fmt.Errorf("no such file: %s, hash: %s %v", audioFile, audioHash, hash)
		return
	}
	if len(action) == 1 {
		var num = -1

		for i, h := range hash {
			if h == audioHash {
				num = i
			}
		}
		if num != -1 {
			var tip model.Tip
			for _, tip = range s.tips {
				if tip.Hash == action[0].StandardWav.Hash[num] {
					break
				}
			}
			if num < len(action[0].StandardWav.Sound) && num != -1 {
				logstash.Info("Input the standard tip audio: " + s.audioRefs[tip.Hash])
				//s.eva.UpdateRefFromVec(path.Join(viper.GetString("AudioRefsDir"), s.audioRefs[tip.Hash]))
			}
		}
	}

	if err = s.audio.Play(audioFile); err != nil {
		return
	}
	if !util.CheckCtx(s.current.ctx) {
		return errs.ErrCtxDone
	}
	return
}

func has(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
