package interact

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"smartconn.cc/tosone/ra-plus/common/logstash"
	"smartconn.cc/tosone/ra-plus/common/online"
	"smartconn.cc/tosone/ra-plus/common/util"
	"smartconn.cc/tosone/ra-plus/drivers/audio"
	"smartconn.cc/tosone/ra-plus/eventCollection"
	"smartconn.cc/tosone/ra-plus/internal/errs"
	"smartconn.cc/tosone/ra-plus/internal/model/event"
	"smartconn.cc/tosone/ra-plus/internal/tips"
	"smartconn.cc/tosone/ra-plus/parsing/lib/betterMatch"
)

// Run 运行某一页的脚本
func (s *State) Run(ctx context.Context, pageID string) (err error) {
	if s.settingFinish != nil {
		s.settingFinish.Wait()
		s.settingFinish = nil
	}

	if err = s.Stop(); err != nil {
		return
	}

	if s.current.ctx != nil && util.CheckCtx(s.current.ctx) && s.current.ctxCancel != nil {
		s.current.ctxCancel()
	}

	s.current = new(Current)
	s.current.ctx, s.current.ctxCancel = context.WithCancel(ctx)
	s.current.pageID = pageID
	s.current.retryTimes = 1
	if s.globalSetting.PerPageQuestion.ProbabilityMode == "num" || s.globalSetting.PerPageQuestion.ProbabilityMode == "" {
		if len(s.questions[s.current.pageID]) == 0 {
			logstash.WithFields(logstash.Fields{"pageID": s.current.pageID}).Error("this page has no interact script")
			err = s.playSyncWithCtx(tips.Get("interact_blank_page"))
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		var randomList []string
		maxQuestion := s.globalSetting.PerPageQuestion.Num // 此模式下最多可以问的问题数目
		if maxQuestion == 0 {
			maxQuestion = 1
		}
		if s.globalSetting.PerPageQuestion.Num > len(s.questions[s.current.pageID]) {
			maxQuestion = len(s.questions[s.current.pageID])
		}
		questionIDList := []string{}
		for str := range s.questions[s.current.pageID] {
			questionIDList = append(questionIDList, str)
		}
		for len(randomList) < maxQuestion {
			rand.Seed(time.Now().UTC().UnixNano())
			if random := util.RandSelect(questionIDList); !hasStr(randomList, random) {
				randomList = append(randomList, random)
			}
		}

		for _, s.current.questionID = range randomList {
			s.current.question = s.questions[s.current.pageID][s.current.questionID]
			if len(s.current.question[0].VeryBefore.Hash) != 0 {
				if e := s.playSyncHashWithCtx(util.RandSelect(s.current.question[0].VeryBefore.Hash)); e == errs.ErrCtxDone {
					err = nil
					return
				} else if e != nil {
					return e
				}
			}
			s.singleQuestion(nil)
		}
	} else if s.globalSetting.PerPageQuestion.ProbabilityMode == "average" {
		for _, q := range s.questions[s.current.pageID] {
			rand.Seed(time.Now().UTC().UnixNano())
			if rand.Intn(len(s.questions[s.current.pageID])) == 0 {
				s.current.question = q
				s.singleQuestion(nil)
			}
		}
	}
	return
}

func (s *State) singleQuestion(flag error) (err error) {
	if flag != nil {
		logstash.WithFields(logstash.Fields{"questionRetry": s.current.retryTimes}).Info(flag)
	}
	if s.current.retryTimes > s.globalSetting.RetryAsk.Num ||
		(flag == errBackWordsWrong && s.globalSetting.BreakBy == "wrong") ||
		(flag == errBackWordsVoid && s.globalSetting.BreakBy == "silence") { // 达到这些条件之后问题就不再重复问了
		logstash.WithFields(logstash.Fields{
			"questionRetry": s.current.retryTimes,
			"RetryAsk.Num":  s.globalSetting.RetryAsk.Num,
			"err":           flag,
		}).Info("cannot get a right response at last")
		if s.current.retryTimes > s.globalSetting.RetryAsk.Num {
			s.neverAgain(flag) // neverAgainWrong or neverAgainNothing
		} else if flag == errBackWordsWrong {
			s.continuousWrong = append(s.continuousWrong, s.current.questionID)
			s.continuousRight = []string{}
			s.continuousSilence = []string{}
			if err = s.ensurePlayLocal("interact_said_wrong", s.current.question[0].SaidWrong.Hash, s.globalSetting.SaidWrong.Hash); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
			if len(s.continuousWrong) >= s.globalSetting.ContinuousBackoutWords.Num {
				if err = s.hookMoreResp("wrong", s.continuousWrong); err != nil {
					logstash.Error(err)
				}
				s.continuousWrong = []string{}
			}
		} else if flag == errBackWordsVoid {
			s.continuousSilence = append(s.continuousSilence, s.current.questionID)
			s.continuousRight = []string{}
			s.continuousWrong = []string{}
			if err = s.ensurePlayLocal("interact_said_nothing", s.current.question[0].SaidNothing.Hash, s.globalSetting.SaidNothing.Hash); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
			if len(s.continuousSilence) >= s.globalSetting.ContinuousBackSilence.Num {
				if err = s.hookMoreResp("silence", s.continuousSilence); err != nil {
					logstash.Error(err)
				}
				s.continuousSilence = []string{}
			}
		}
		return
	}
	switch flag {
	case errBackWordsVoid:
		if err = s.ensurePlayLocal("interact_said_nothing", s.current.question[0].SaidNothing.Hash, s.globalSetting.SaidNothing.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
	case errBackWordsWrong:
		if err = s.ensurePlayLocal("interact_said_wrong", s.current.question[0].SaidWrong.Hash, s.globalSetting.SaidWrong.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
	}
	if !util.CheckCtx(s.current.ctx) {
		return
	}
	for _, q := range s.current.question {
		s.current.keywords = append(s.current.keywords, q.Keyword...)
	}
	if len(s.current.question) > 0 && len(s.current.question[0].Before.Hash) != 0 {
		var questionTip string
		if s.current.retryTimes != 1 && len(s.current.question[0].AgainBefore.Hash) != 0 { // 取出第二次问题 hash
			questionTip = util.RandSelect(s.current.question[0].AgainBefore.Hash)
		} else {
			questionTip = util.RandSelect(s.current.question[0].Before.Hash) // 取出问题 hash
		}
		if err = s.playSyncHashWithCtx(questionTip); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		if s.current.question[0].NoRecord { // 不需要录音的时候所做的操作
			logstash.Info("question is over with noRecord args set true")
		} else {
			s.questionRecord() // 开始录音
		}
	} else { // 播放这一页什么都没有
		if err = s.playSyncWithCtx(tips.Get("interact_blank_page")); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
	}
	return
}

// questionRecord ..
func (s *State) questionRecord() (err error) {
	var VoiceDetCallback, STTResultCallback, STTBadReportCallback func()
	var TimerVoiceDetectTimeout *time.Timer

	s.current.sttResult = "" // 每次录音循环中初始化

	recordStartStamp := time.Now()

	channel := make(chan error, 1)

	logstash.Info("new question record entered")

	var ret = func() {
		if channel != nil {
			close(channel)
			channel = nil
		}
		audio.RecordStop()
		s.vi.Kill(false)
		if STTBadReportCallback != nil {
			STTBadReportCallback()
		}
		if VoiceDetCallback != nil {
			VoiceDetCallback()
		}
		if STTResultCallback != nil {
			STTResultCallback()
		}
		if TimerVoiceDetectTimeout != nil {
			TimerVoiceDetectTimeout.Stop()
		}
	}

	defer ret()

	if !util.CheckCtx(s.current.ctx) {
		return
	}

	s.vi.WaitUntilReady()

	if !util.CheckCtx(s.current.ctx) {
		return
	}

	logstash.WithFields(logstash.Fields{"loop": s.current.retryTimes}).Info("start question loop")

	var reader io.Reader
	reader, err = audio.RecordStart()
	if err != nil {
		logstash.Error(err.Error())
		return
	}

	logstash.Info("record is starting")

	TimerVoiceDetectTimeout = time.AfterFunc(time.Second*time.Duration(s.globalSetting.MaxRecording), func() { // 到达最长的录音时间
		logstash.Info("at the end said nothing")
		if !util.CheckCtx(s.current.ctx) {
			return
		}
		channel <- errVoiceDetectTimeout
	})

	STTBadReportCallback = s.vi.OnNetSttBadReport(func() {
		if !util.CheckCtx(s.current.ctx) {
			return
		}
		channel <- errBadNetwork
	})

	STTResultCallback = s.vi.OnSttReport(func(str string) {
		s.current.sttResult = str
		logstash.WithFields(logstash.Fields{"sttResult": str}).Info("voice detect with stt result string")
		if !util.CheckCtx(s.current.ctx) {
			return
		}
		channel <- errVoiceDetectSuccess
	})

	VoiceDetCallback = s.vi.OnVoiceDet(func(flag bool) {
		logstash.WithFields(logstash.Fields{"flag": flag}).Info("voice detect with a flag")
		if !util.CheckCtx(s.current.ctx) {
			return
		}
		if !flag { // 说话声结束
			audio.RecordStop()
		}
	})

	if !util.CheckCtx(s.current.ctx) {
		return
	}

	if !online.Get() {
		if err = s.playSyncWithCtx(tips.Get("interact_offline")); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		return
	}

	if err = s.vi.Start(reader); err != nil {
		if err = s.playSyncWithCtx(tips.Get("interact_offline")); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		} // 播放网络状态不好
		return
	}

	info := <-channel

	logstash.Info(fmt.Sprintf("current question loop is over with code: %v", info))

	s.current.responseTime = time.Since(recordStartStamp).Seconds()

	ret()

	if info == nil {
		return
	}

	logstash.Info(fmt.Sprintf("question loop is over with code: %v", info))
	if s.current.sttResult == "" {
		logstash.Info("stt nothing response, wait for a second to get result")
		<-time.After(time.Millisecond * 1000)
		if !util.CheckCtx(s.current.ctx) {
			return
		}
		logstash.WithFields(logstash.Fields{"stt": s.current.sttResult}).Info("one second later")
	}
	logstash.Info("clear all of the status")
	if strings.HasPrefix(s.current.sttResult, "\"") {
		s.current.sttResult = strings.Trim(s.current.sttResult, "\"")
	}
	s.current.retryTimes++
	switch info {
	case errVoiceDetectSuccess:
		if err := s.questionResp(s.current.sttResult); err != nil {
			s.singleQuestion(err)
		}
	case errVoiceDetectTimeout:
		if s.current.sttResult != "" {
			if err := s.questionResp(s.current.sttResult); err != nil {
				s.singleQuestion(err)
			}
		} else {
			s.hookResp()
			s.singleQuestion(errBackWordsVoid)
		}
	case errBadNetwork: // 网络差导致上传音频数据到语音识别 server 太慢产生的错误
		if s.current.sttResult != "" {
			if err := s.questionResp(s.current.sttResult); err != nil {
				s.singleQuestion(err)
			}
		} else {
			s.singleQuestion(errBackWordsVoid)
		}
	}
	return
}

// questionResp 处理正确回答部分的内容
func (s *State) questionResp(sttResult string) (err error) {
	if sttResult == "" {
		s.hookResp()
		return errBackWordsVoid
	}
	var result string
	if result = betterMatch.Match(sttResult, s.current.keywords); result == "" {
		s.hookResp()
		return errBackWordsWrong
	}
	for _, q := range s.current.question {
		if hasStr(q.Keyword, result) {
			logstash.WithFields(logstash.Fields{"keywords": q.Keyword, "sttResult": result}).Info("user said something got a right response")
			s.continuousRight = append(s.continuousRight, s.current.questionID)
			s.continuousWrong = []string{}
			s.continuousSilence = []string{}
			logstash.WithFields(logstash.Fields{"continueRight": len(s.continuousRight)}).Info("continue count")

			s.hookResp(q.Keyword...)

			if err = s.ensurePlayLocal("interact_said_right", q.Hash); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}

			if len(q.Award.Hash) != 0 { // 播放奖励音频
				awardSound := util.RandSelect(q.Award.Hash)
				if err = s.playSyncHashWithCtx(awardSound); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			}

			if len(s.continuousRight) >= s.globalSetting.ContinuousBackinWords.Num {
				err = s.hookMoreResp("right", s.continuousRight)
				if err = s.ensurePlayLocal("interact_continuous_backin_words", s.globalSetting.ContinuousBackinWords.Hash); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			}
			return
		}
	}

	return
}

func (s *State) neverAgain(flag error) (err error) {
	var gotAudio bool
	switch flag {
	case errBackWordsVoid:
		if err = s.ensurePlayLocal("interact_never_again_silence", s.current.question[0].NeverAgainNothing.Hash, s.globalSetting.NeverAgainNothing.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		gotAudio = true
	case errBackWordsWrong:
		if err = s.ensurePlayLocal("interact_never_again_wrong", s.current.question[0].NeverAgainWrong.Hash, s.globalSetting.NeverAgainWrong.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		gotAudio = true
	default:
		logstash.Info(fmt.Sprintf("neverAgain: no suck a command %s.", flag))
	}
	if !gotAudio {
		if err = s.ensurePlayLocal("interact_never_again", s.current.question[0].NeverAgain.Hash, s.globalSetting.NeverAgain.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
	}
	return
}

func (s *State) hookMoreResp(t string, questions []string) error {
	return eventCollection.Add("interact", event.MoreResp{
		Event:      "continuous-answer",
		BookID:     s.bookID,
		PageID:     s.current.pageID,
		UniqueID:   s.uniqueID,
		QuestionID: s.current.questionID,
		Time:       time.Now().Unix() * 1000,
		Count:      s.globalSetting.ContinuousBackSilence.Num,
		Questions:  questions,
		Type:       t,
		Mode:       "interact",
	})
}

func (s *State) hookResp(key ...string) error {
	return eventCollection.Add("parsing", event.QuestionResp{
		Event:        "answer-question",
		BookID:       s.bookID,
		PageID:       s.current.pageID,
		UniqueID:     s.uniqueID,
		Time:         time.Now().Unix() * 1000,
		QuestionID:   s.current.questionID,
		ResponseTime: 1,
		AnswerTimes:  s.current.retryTimes,
		Award:        s.current.award,
		STTResult:    s.current.sttResult,
		Mode:         "interact",
		KeywordList:  key,
	}) // 事件上报
}
