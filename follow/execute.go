package follow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"smartconn.cc/tosone/parsing-go/errs"
	"smartconn.cc/tosone/parsing-go/logstash"
	"smartconn.cc/tosone/parsing-go/model/event"
	"smartconn.cc/tosone/parsing-go/tips"
	"smartconn.cc/tosone/parsing-go/util"
)

func (s *State) Run(ctx context.Context, pageID string) (err error) {
	if s.settingFinish != nil {
		s.settingFinish.Wait()
		s.settingFinish = nil
	}

	if err = s.Stop(); err != nil {
		return
	}

	if !util.CheckCtx(ctx) {
		return
	}

	s.current.ctx, s.current.ctxCancel = context.WithCancel(ctx)
	s.current.pageID = pageID

	if pageID == s.bookID {
		return s.Intro()
	}

	return s.Event(pageID)

	return
}

// Intro 任务介绍
func (s *State) Intro() (err error) {
	logstash.Info("Playing the intro task tip audio.")

	if s.mixedSetting.IntroTask == nil {
		return errors.New("intro task audio not set yet")
	}

	if len(*s.mixedSetting.IntroTask) == 0 {
		return errors.New("intro task audio length not correct")
	}

	length := len((*s.mixedSetting.IntroTask)[0].Sound)

	if length != 2 {
		return errors.New("intro task audio length not correct")
	}

	for _, item := range *s.mixedSetting.IntroTask {
		if length != len(item.Sound) {
			return errors.New("intro task audio length not correct")
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	selected := rand.Intn(len((*s.mixedSetting.IntroTask)))

	if err = s.playWithCtx([]string{(*s.mixedSetting.IntroTask)[selected].Sound[0]}, []string{(*s.mixedSetting.IntroTask)[0].Hash[selected]}); err != nil {
		if err == errs.ErrCtxDone {
			err = nil
		}
		return
	}

	if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.totalFlower))); err != nil {
		if err == errs.ErrCtxDone {
			err = nil
		}
		return
	}

	if err = s.playWithCtx([]string{(*s.mixedSetting.IntroTask)[selected].Sound[1]}, []string{(*s.mixedSetting.IntroTask)[selected].Hash[1]}); err != nil {
		if err == errs.ErrCtxDone {
			err = nil
		}
		return
	}
	return
}

func (s *State) calcTotalFlower() (err error) {
	strList := strings.Split(s.mixedSetting.Goal.FlowerPerTask, "/")
	if len(strList) != 2 {
		logstash.Error(fmt.Sprintf("%s have a syntx error: FlowerPerTask %s not expected.",
			s.mixedSetting.Goal.FlowerPerTask, s.mixedSetting.Goal.FlowerPerTask))
		return
	}

	var flowerNum int
	var taskNum int

	if flowerNum, err = strconv.Atoi(strList[0]); err != nil {
		return
	}
	if taskNum, err = strconv.Atoi(strList[1]); err != nil {
		return
	}
	s.totalFlower = (flowerNum / taskNum) * s.mixedSetting.Goal.Task
	return
}

// Takeoff 拔下书
func (s *State) Takeoff() {
	var err error
	s.Stop()
	if s.mixedSetting.TakeOff != nil {
		if s.goal {
			if err = s.playWithCtx(s.mixedSetting.TakeOff.Success.Sound, s.mixedSetting.TakeOff.Success.Hash); err != nil {
				logstash.Error(err.Error())
			}
		} else {
			if err = s.playWithCtx(s.mixedSetting.TakeOff.Fail.Sound, s.mixedSetting.TakeOff.Fail.Hash); err != nil {
				logstash.Error(err.Error())
			}
		}
	}

	//eventinteraction.EventReportAdd("followReading", event.FollowRecordTakeOff{
	//	Event:           "follow-reading-end",
	//	InsertID:        s.insertID,
	//	BookID:          s.bookID,
	//	UniqueID:        s.uniqueID,
	//	Task:            uint(s.mixedSetting.Goal.Task),
	//	FinishTask:      uint(s.task),
	//	ExtraTask:       uint(s.mixedSetting.Extra.Task),
	//	FinishExtraTask: uint(s.extraTask),
	//	Flower:          uint(s.flower),
	//	ValidFollow:     s.validFollow,
	//})
}

// Event 开始跟读，需要传入当前的页面 PageID，然后开始这一页的跟读
func (s *State) Event(pageID string) (err error) {
	//var reader io.Reader

	if s.nextPageTimer != nil {
		s.nextPageTimer.Stop()
	}
	s.Stop() // stop the last follow flow
	s.ctxWg = new(sync.WaitGroup)
	s.ctxWg.Add(1)
	defer func() {
		s.ctxWg.Done()
		s.ctxWg = nil
	}()
	//s.current.ctx, s.current.ctxCancel = context.WithCancel(ctx)
	defer s.autoTipNextPage() // 结束的时候自动开始自动倒计时提示下一页的提示音频

	var script FollowItem
	var ok bool
	logstash.Info(fmt.Sprintf("follow scripts: %+v", s.allPages))
	if script, ok = s.allPages[pageID]; !ok {
		if err = s.playPureWithCtx(tips.Get("interact_blank_page")); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
		return
	}
	if !s.contentEntered && s.mixedSetting.ContentBefore != nil {
		if !has(pageID, s.mixedSetting.ExcludeContent) {
			s.contentEntered = true
			logstash.WithFields(logstash.Fields{
				"exclude": s.mixedSetting.ExcludeContent,
				"pageID":  pageID,
			}).Info("Playing contentBefore tip audio.")
			if err = s.playWithCtx(s.mixedSetting.ContentBefore.Sound,
				s.mixedSetting.ContentBefore.Hash); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
		}
	}
	if script.NoRecord {
		logstash.Info("This page is null.")
		if err = s.playPureWithCtx(tips.Get("interact_blank_page")); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
		}
		return
	}
	logstash.Info(fmt.Sprintf("follow scripts: %+v", script.Follows))

	for _, followID := range script.Order {
		var entryList []string
		var evaMode string // last,average default: last
		if _, ok := script.InnerOrder[followID]; ok {
			entryList = script.InnerOrder[followID]
			if action, ok := script.Follows[entryList[len(entryList)-1]]; ok {
				for _, item := range action.StandardOgg.Sound {
					if len(strings.Split(strings.TrimPrefix(item, "oss://bbcloud-story/"), "_")) == 3 {
						evaMode = "last"
					} else {
						evaMode = "average"
					}
				}
			}
		} else {
			entryList = []string{followID}
			evaMode = "last"
		}

		logstash.Info(fmt.Sprintf("Entry follow list: %+v, mode: %s", entryList, evaMode))

		var averageScore float64
		var recordFileList []string
		var taskContentList [][]byte

		for entryIndex, entry := range entryList {
			if action, ok := script.Follows[entry]; !ok {
				logstash.Error(fmt.Sprintf("Script has no followID: %s", entry))
			} else {
				logstash.Info(fmt.Sprintf("Follow id: %s", entry))
				var num uint
				var score uint // 录音的最终得分
				for {
					select {
					case <-s.current.ctx.Done():
						return
					default:
					}
					num++

					if num > s.mixedSetting.Retry.MaxRetry {
						logstash.Info("Got its max retry times.")
						break
					}
					logstash.Info("Playing standard tip audio.")
					if err = s.playWithCtx(action.StandardOgg.Sound, action.StandardOgg.Hash, action); err != nil {
						if err == errs.ErrCtxDone {
							err = nil
						}
						return
					}
					if s.mixedSetting.BeforeRecord != nil {
						if err = s.playWithCtx(s.mixedSetting.BeforeRecord.Sound,
							s.mixedSetting.BeforeRecord.Hash); err != nil {
							if err == errs.ErrCtxDone {
								err = nil
							}
							return
						}
					} else {
						if err = s.audio.Play(tips.Get("scanning_page_di")); err != nil {
							if err == errs.ErrCtxDone {
								err = nil
							}
							return
						}
					}

					//var recordOver = make(chan bool, 1) // 录音是否已经结束
					//var maxRecordingTimer *time.Timer   // 最长录音时间的定时器
					//ledutil.Flicker()
					//var gotVoiceEdge <-chan bool
					//var recordRandFile = util.UUID()
					var followStart = time.Now().Unix() * 1000
					var followEnd int64
					var recordFile string
					s.audio.Stop()
					//if reader, err = audio.RecordStart(); err != nil { // 暂时不能录音
					//	logstash.Error(err.Error())
					//} else {
					//	var fileReader, fileWriter = io.Pipe()
					//	var evaDetReader, evaDetWriter = io.Pipe()
					//	mw := io.MultiWriter(fileWriter, evaDetWriter)
					//	go func() {
					//		if _, err = io.Copy(mw, reader); err != nil {
					//			logstash.Error(err.Error())
					//		}
					//		fileReader.CloseWithError(io.EOF)
					//		evaDetReader.CloseWithError(io.EOF)
					//	}()
					//	var recordFileHandle *os.File
					//	recordFile = path.Join(viper.GetString("FollowRecordDir"), recordRandFile)
					//	if recordFileHandle, err = os.Create(recordFile); err != nil {
					//		logstash.Error(err.Error())
					//	}
					//
					//	go func() {
					//		for {
					//			var data = make([]byte, 1024)
					//			var err error
					//			var n int
					//			n, err = fileReader.Read(data)
					//			if err == io.EOF {
					//				break
					//			}
					//			if err != nil {
					//				logstash.Info(err.Error())
					//				break
					//			}
					//			if n == 0 {
					//				continue
					//			}
					//			if _, err = recordFileHandle.Write(data[:n]); err != nil {
					//				logstash.Info(err.Error())
					//				break
					//			}
					//		}
					//	}()
					//	gotVoiceEdge = s.eva.WriteAudio(evaDetReader)
					//}
					//
					//maxRecordingTimer = time.AfterFunc(time.Second*time.Duration(*s.mixedSetting.MaxRecording), func() {
					//	recordOver <- true
					//})
					//select {
					//case <-gotVoiceEdge:
					//	logstash.Info("Got voice edge. Stop record.")
					//case <-recordOver:
					//	logstash.Info("Timeout to record")
					//case <-s.current.ctx.Done():
					//	if maxRecordingTimer != nil {
					//		maxRecordingTimer.Stop()
					//	}
					//	return
					//}
					//followEnd = time.Now().Unix() * 1000

					s.audio.Stop()
					//var correctScore, speedScore, toneScore uint
					//score, correctScore, toneScore, speedScore = s.eva.GetEvaResult()
					//if maxRecordingTimer != nil {
					//	maxRecordingTimer.Stop()
					//}
					//score = 90 // need remove
					logstash.Info(fmt.Sprintf("Follow record got score: %d", score))
					logstash.Info(fmt.Sprintf("Follow record raw audio: %s.", recordFile))
					var taskContent []byte
					if taskContent, err = json.Marshal(event.FollowRecordInfo{
						Event:    "follow-reading",
						InsertID: s.insertID,
						BookID:   s.bookID,
						UniqueID: s.uniqueID,
						PageID:   pageID,
						FollowID: entry,
						Start:    followStart,
						End:      followEnd,
						Score:    score,
						Detail: event.FollowRecordInfoDetail{
							Correct: 10,
							Speed:   10,
							Tone:    10,
						},
						Times:  num,
						Record: recordFile,
					}); err != nil {
						logstash.Error(err.Error())
					} else {
						taskContentList = append(taskContentList, taskContent)
					}

					if score != 0 && (num == s.mixedSetting.Retry.MaxRetry || score > s.mixedSetting.Retry.Score) {
						recordFileList = append(recordFileList, recordFile)
					}

					if entryIndex == len(entryList)-1 && (num == s.mixedSetting.Retry.MaxRetry || score > s.mixedSetting.Retry.Score) {
						if evaMode == "last" {
							if err = s.playPCMWithCtx(recordFile); err != nil {
								if err == errs.ErrCtxDone {
									err = nil
								}
								return
							}
							if err = s.playWithCtx(action.StandardOgg.Sound, action.StandardOgg.Hash); err != nil {
								if err == errs.ErrCtxDone {
									err = nil
								}
								return
							}
						} else {
							for _, item := range recordFileList {
								if err = s.playPCMWithCtx(item); err != nil {
									if err == errs.ErrCtxDone {
										err = nil
									}
									return
								}
							}
							for _, entry := range entryList {
								if action, ok := script.Follows[entry]; !ok {
									logstash.Error(fmt.Sprintf("Script has no followID: %s", entry))
								} else {
									if err = s.playWithCtx(action.StandardOgg.Sound, action.StandardOgg.Hash); err != nil {
										if err == errs.ErrCtxDone {
											err = nil
										}
										return
									}
								}
							}
						}
					}

					if num != s.mixedSetting.Retry.MaxRetry {
						if score == 0 {
							if err = s.scoreCalc(followID, score, entryIndex == len(entryList)-1); err != nil {
								logstash.Error(err.Error()) // got this err just continue
							}
						} else if score <= s.mixedSetting.Retry.Score {
							if err = s.playWithCtx(
								s.mixedSetting.Retry.BadScore.Sound,
								s.mixedSetting.Retry.BadScore.Hash); err != nil {
								if err == errs.ErrCtxDone {
									err = nil
								}
								return
							}
						}
					}

					if score > s.mixedSetting.Retry.Score {
						logstash.Info("Got a score above the minimum line.")
						if averageScore == 0 {
							averageScore = float64(score)
						} else {
							averageScore = (averageScore + float64(score)) / 2
						}
						break
					} else if num == s.mixedSetting.Retry.MaxRetry {
						if averageScore == 0 {
							averageScore = float64(score)
						} else {
							averageScore = (averageScore + float64(score)) / 2
						}
					}
				}
			}
		}

		if err = s.scoreCalc(followID, uint(averageScore), true); err != nil {
			logstash.Error(err.Error()) // got this err just continue
		}

		//for _, item := range taskContentList {
		//	//store.AddTask("followRecordUpload", item)
		//}
	}

	if pageID == s.mixedSetting.LastContentPage {
		if s.mixedSetting.ContentAfterSuccess != nil && s.mixedSetting.ContentAfterFail != nil {
			if s.goal {
				logstash.Info("Playing content after success tip audio.")
				if err = s.playWithCtx(s.mixedSetting.ContentAfterSuccess.Sound,
					s.mixedSetting.ContentAfterSuccess.Hash); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			} else {
				if s.mixedSetting.ContentAfterFail == nil {
					return errors.New("content after tip audio not set yet")
				}

				if len(*s.mixedSetting.ContentAfterFail) == 0 {
					return errors.New("content after tip audio length not correct")
				}

				length := len((*s.mixedSetting.ContentAfterFail)[0].Sound)

				if length != 3 {
					return errors.New("content after tip audio length not correct")
				}

				for _, item := range *s.mixedSetting.ContentAfterFail {
					if length != len(item.Sound) {
						return errors.New("content after tip audio length not correct")
					}
				}

				rand.Seed(time.Now().UTC().UnixNano())
				selected := rand.Intn(len((*s.mixedSetting.ContentAfterFail)))
				logstash.Info("Playing content after fail tip audio.")
				if err = s.playWithCtx([]string{(*s.mixedSetting.ContentAfterFail)[selected].Sound[0]},
					[]string{(*s.mixedSetting.ContentAfterFail)[selected].Hash[0]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.flower))); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playWithCtx([]string{(*s.mixedSetting.ContentAfterFail)[selected].Sound[1]},
					[]string{(*s.mixedSetting.ContentAfterFail)[selected].Hash[1]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.totalFlower-s.flower))); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playWithCtx([]string{(*s.mixedSetting.ContentAfterFail)[selected].Sound[2]},
					[]string{(*s.mixedSetting.ContentAfterFail)[selected].Hash[2]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			}
		}
	} else {
		logstash.Info("Playing next page tip audio.")
		if s.mixedSetting.NextPage != nil {
			if err = s.playWithCtx(s.mixedSetting.NextPage.Sound,
				s.mixedSetting.NextPage.Hash); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
		}
	}

	return
}

func (s *State) autoTipNextPage() {
	var err error
	s.nextPageTimer = time.AfterFunc(time.Second*60, func() {
		if err = s.playWithCtx(s.mixedSetting.NextPage.Sound, s.mixedSetting.NextPage.Hash); err != nil {
			if err == errs.ErrCtxDone {
				err = nil
			}
			return
		}
	})
}

// maxFollowInvalid 判断当前的跟读是否已经达到当前的跟读最大次数
func (s *State) maxFollowInvalid(sentenceID string) bool {
	for id, times := range s.sentenceMap {
		if id == sentenceID {
			if times >= 3 {
				return true // 再次得分无效
			}
		}
	}
	return false
}

// flower 判断当前的得分可以获得多少朵小红花，是否算是完整这个跟读，是否达到了某个进度
func (s *State) scoreCalc(sentenceID string, score uint, last bool) (err error) {
	var evaluation Evaluation

	if last {
		if score != 0 { // 得分无效之后直接提示得分无效，不再评价用户的语音
			if num, ok := s.sentenceMap[sentenceID]; ok {
				if num < s.mixedSetting.MaxFollowInvalid.Num {
					s.sentenceMap[sentenceID]++
				} else { // 得分已经无效
					logstash.Info(fmt.Sprintf("Playing max follow invalid sentenceID: %s, times: %d.", sentenceID, s.sentenceMap[sentenceID]))
					if err = s.playWithCtx(s.mixedSetting.MaxFollowInvalid.Sound, s.mixedSetting.MaxFollowInvalid.Hash); err != nil {
						if err == errs.ErrCtxDone {
							err = nil
						}
						return
					}
				}
			} else {
				s.sentenceMap[sentenceID] = 1
			}
		}
	}

	var valid bool // 跟读是否有效
	if s.mixedSetting.Evaluation != nil {
		for _, evaluation = range *s.mixedSetting.Evaluation {
			if score == 0 && evaluation.MaxScore == 0 && evaluation.MinScore == 0 {
				valid = evaluation.Valid
				logstash.Info("Playing evaluation 0 tip audio.")
				if err = s.playWithCtx(evaluation.Sound, evaluation.Hash); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				break
			} else if score >= evaluation.MinScore && score < evaluation.MaxScore && last {
				valid = evaluation.Valid
				logstash.Info(fmt.Sprintf("Playing evaluation max: %d, min: %d tip audio.", evaluation.MaxScore, evaluation.MinScore))
				if err = s.playWithCtx(evaluation.Sound, evaluation.Hash); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				break
			}
		}
	}

	if !valid || !last { // 得分无效或者不是最后一句
		return
	}

	s.validFollow++
	var sentenceDoneNum int // 已经完成的句子的数量，句子重复得分也是有效的，但有最大有效次数
	for _, val := range s.sentenceMap {
		sentenceDoneNum += val
	}

	// 计算小红花是否可以获得
	if !s.goal { // 没有完成目标之前
		strs := strings.Split(s.mixedSetting.Goal.SentencePerFlower, "/")
		if len(strs) != 2 {
			logstash.Error(fmt.Sprintf("%s have a syntx error: FlowerPerTask %s not expected.", sentenceID, s.mixedSetting.Goal.SentencePerFlower))
			return
		}
		var sentenceNum int
		var flowerNum int
		if sentenceNum, err = strconv.Atoi(strs[0]); err != nil {
			return
		}
		if flowerNum, err = strconv.Atoi(strs[1]); err != nil {
			return
		}
		//sentenceNum = 1 // need to remove
		//flowerNum = 1   // need to remove
		if sentenceDoneNum%sentenceNum == 0 {
			s.flower += flowerNum
		} else { // 暂时无法获得小红花
			return
		}
		logstash.WithFields(logstash.Fields{
			"sentenceNum": sentenceNum,
			"flowerNum":   flowerNum, "Flower": s.flower,
			"SentencePerFlower": s.mixedSetting.Goal.SentencePerFlower,
		}).Info("Current flower status.")
	} else {
		strs := strings.Split(s.mixedSetting.Extra.SentencePerFlower, "/")
		if len(strs) != 2 {
			logstash.Error(fmt.Sprintf("%s have a syntx error: FlowerPerTask %s not expected.",
				sentenceID, s.mixedSetting.Extra.SentencePerFlower))
			return
		}
		var sentenceNum int
		var flowerNum int
		if sentenceNum, err = strconv.Atoi(strs[0]); err != nil {
			return
		}
		if flowerNum, err = strconv.Atoi(strs[1]); err != nil {
			return
		}
		if sentenceDoneNum%sentenceNum == 0 {
			s.extraFlower += flowerNum
		} else { // 暂时无法获得小红花
			return
		}
	}

	if !s.goal { // 没有完成目标之前
		strs := strings.Split(s.mixedSetting.Goal.FlowerPerTask, "/")
		if len(strs) != 2 {
			logstash.Error(fmt.Sprintf("%s have a syntx error: FlowerPerTask %s not expected.",
				sentenceID, s.mixedSetting.Goal.FlowerPerTask))
			return
		}

		var flowerNum int
		var taskNum int

		if flowerNum, err = strconv.Atoi(strs[0]); err != nil {
			return
		}
		if taskNum, err = strconv.Atoi(strs[1]); err != nil {
			return
		}
		s.totalFlower = (flowerNum / taskNum) * s.mixedSetting.Goal.Task
		//flowerNum = 1 // need to remove
		//taskNum = 1   // need to remove
		if s.flower%flowerNum == 0 { // 达到一个新的进度
			s.task += taskNum
			logstash.Info(fmt.Sprintf("Got flower %d", s.flower))
			if s.mixedSetting.Progress == nil || len(*s.mixedSetting.Progress) != 3 {
				return
			}
			if s.task < s.mixedSetting.Goal.Task {
				if s.mixedSetting.Progress == nil {
					return errors.New("progress tip audio not set yet")
				}

				if len(*s.mixedSetting.Progress) == 0 {
					return errors.New("progress tip audio length not correct")
				}

				length := len((*s.mixedSetting.Progress)[0].Sound)

				if length != 3 {
					return errors.New("progress tip audio length not correct")
				}

				for _, item := range *s.mixedSetting.Progress {
					if length != len(item.Sound) {
						return errors.New("progress tip audio length not correct")
					}
				}

				rand.Seed(time.Now().UTC().UnixNano())
				selected := rand.Intn(len((*s.mixedSetting.Progress)))
				logstash.Info(fmt.Sprintf("progress tip audio select: %d", selected))

				if err = s.playWithCtx([]string{(*s.mixedSetting.Progress)[selected].Sound[0]},
					[]string{(*s.mixedSetting.Progress)[selected].Hash[0]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.flower))); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playWithCtx([]string{(*s.mixedSetting.Progress)[selected].Sound[1]},
					[]string{(*s.mixedSetting.Progress)[selected].Hash[1]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.totalFlower-s.flower))); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
				if err = s.playWithCtx([]string{(*s.mixedSetting.Progress)[selected].Sound[2]},
					[]string{(*s.mixedSetting.Progress)[selected].Hash[2]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			} else if s.task == s.mixedSetting.Goal.Task {
				s.goal = true

				if s.mixedSetting.ProgressComplete == nil {
					return errors.New("progress complete tip audio not set yet")
				}

				if len(*s.mixedSetting.ProgressComplete) == 0 {
					return errors.New("progress complete tip audio length not correct")
				}

				length := len((*s.mixedSetting.ProgressComplete)[0].Sound)

				if length != 2 {
					return errors.New("progress complete tip audio length not correct")
				}

				for _, item := range *s.mixedSetting.ProgressComplete {
					if length != len(item.Sound) {
						return errors.New("progress complete tip audio length not correct")
					}
				}

				rand.Seed(time.Now().UTC().UnixNano())
				selected := rand.Intn(len((*s.mixedSetting.ProgressComplete)))

				if err = s.playWithCtx([]string{(*s.mixedSetting.ProgressComplete)[selected].Sound[0]},
					[]string{(*s.mixedSetting.ProgressComplete)[selected].Hash[0]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}

				if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.flower))); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}

				if err = s.playWithCtx([]string{(*s.mixedSetting.ProgressComplete)[selected].Sound[1]},
					[]string{(*s.mixedSetting.ProgressComplete)[selected].Hash[1]}); err != nil {
					if err == errs.ErrCtxDone {
						err = nil
					}
					return
				}
			}
		} else { // 暂时无法达到一个任务进度
			return
		}
	} else {
		strs := strings.Split(s.mixedSetting.Extra.FlowerPerTask, "/")
		if len(strs) != 2 {
			logstash.Error(fmt.Sprintf("%s have a syntx error: FlowerPerTask %s not expected.",
				sentenceID, s.mixedSetting.Extra.FlowerPerTask))
			return
		}

		var flowerNum int
		var taskNum int

		if flowerNum, err = strconv.Atoi(strs[0]); err != nil {
			return
		}
		if taskNum, err = strconv.Atoi(strs[1]); err != nil {
			return
		}
		if s.extraFlower%flowerNum == 0 {
			s.extraTask += taskNum

			logstash.Info(fmt.Sprintf("Extra got flower %d", s.extraFlower))
			if s.mixedSetting.ExtraProgress == nil {
				return errors.New("extra progress tip audio not set yet")
			}

			if len(*s.mixedSetting.ExtraProgress) == 0 {
				return errors.New("extra progress complete tip audio length not correct")
			}

			length := len((*s.mixedSetting.ExtraProgress)[0].Sound)

			if length != 2 {
				return errors.New("extra progress complete tip audio length not correct")
			}

			for _, item := range *s.mixedSetting.ExtraProgress {
				if length != len(item.Sound) {
					return errors.New("extra progress complete tip audio length not correct")
				}
			}

			rand.Seed(time.Now().UTC().UnixNano())
			selected := rand.Intn(len((*s.mixedSetting.ExtraProgress)))

			if err = s.playWithCtx([]string{(*s.mixedSetting.ExtraProgress)[selected].Sound[0]},
				[]string{(*s.mixedSetting.ExtraProgress)[selected].Hash[0]}); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
			if err = s.playPureWithCtx(tips.Get(fmt.Sprintf("flower_%d", s.extraFlower))); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
			if err = s.playWithCtx([]string{(*s.mixedSetting.ExtraProgress)[selected].Sound[1]},
				[]string{(*s.mixedSetting.ExtraProgress)[selected].Hash[1]}); err != nil {
				if err == errs.ErrCtxDone {
					err = nil
				}
				return
			}
		} else { // 暂时无法达到一个任务进度
			return
		}
	}

	return
}

// Stop 停止跟读模式
//func (inst *Inst) Stop() {
//	if inst.ctxCancel != nil {
//		inst.ctxCancel()
//		if inst.ctxWg != nil {
//			inst.ctxWg.Wait()
//		}
//	}
//	audioCtrl.Clear()
//	ledutil.Close()
//}

// IsStopped 查询是否已经停止，true 代表已经停止
//func (inst *Inst) IsStopped() bool {
//	select {
//	case <-inst.ctx.Done():
//		return true
//	default:
//		return false
//	}
//	return false
//}

//func has(item string, list []string) bool {
//	for _, i := range list {
//		if i == item {
//			return true
//		}
//	}
//	return false
//}
