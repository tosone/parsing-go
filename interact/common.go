package interact

import (
	"fmt"

	"github.com/Unknwon/com"
	"smartconn.cc/tosone/ra-plus/common/logstash"
	"smartconn.cc/tosone/ra-plus/common/util"
	"smartconn.cc/tosone/ra-plus/drivers/audio"
	"smartconn.cc/tosone/ra-plus/internal/errs"
	"smartconn.cc/tosone/ra-plus/internal/tips"
)

func (s *State) playSyncWithCtx(p string) (err error) {
	if !util.CheckCtx(s.current.ctx) {
		err = errs.ErrCtxDone
		return
	}
	if err = audio.PlaySync(p); err != nil {
		return
	}
	if !util.CheckCtx(s.current.ctx) {
		err = errs.ErrCtxDone
		return
	}
	return
}

func (s *State) playSyncHashWithCtx(p string) (err error) {
	if !util.CheckCtx(s.current.ctx) {
		err = errs.ErrCtxDone
		return
	}

	var file string
	if file = s.tipFromHash(p); file == "" {
		err = fmt.Errorf("no such a file hash: %s", p)
	}

	if !util.CheckCtx(s.current.ctx) {
		err = errs.ErrCtxDone
		return
	}

	err = s.playSyncWithCtx(file)
	return
}

func (s *State) tipFromHash(hash string) (tip string) {
	for _, v := range s.tips {
		if v.Hash == hash {
			tip = v.Path
			return
		}
	}
	return
}

// ensureTip tips, hash, hash
func (s *State) ensureTip(list ...[]string) (tip string) {
	if len(list) == 0 {
		return
	} else if len(list) == 1 {
		tip = util.RandSelect(list[0])
	} else {
		for i := 1; i < len(list); i++ {
			f := s.tipFromHash(util.RandSelect(list[i]))
			if com.IsFile(f) {
				tip = f
				break
			}
		}
		if tip == "" {
			tip = util.RandSelect(list[0])
		}
	}
	return
}

func (s *State) ensurePlay(list ...[]string) (err error) {
	if t := s.ensureTip(list...); t == "" {
		err = fmt.Errorf("cannot find any tip audio")
		return
	} else {
		err = s.playSyncWithCtx(t)
	}
	return
}

func (s *State) ensurePlayLocal(str string, list ...[]string) (err error) {
	var ts []string
	if ts, err = tips.Gets(str); err != nil {
		logstash.Error(err)
		err = nil
	}
	err = s.ensurePlay(ts, s.current.question[0].SaidNothing.Hash, s.globalSetting.SaidNothing.Hash)
	return
}

func hasStr(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
