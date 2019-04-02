package interact

import (
	"context"
	"sync"

	"smartconn.cc/tosone/ra-plus/drivers/voiceinteract/voiceInteract"
	"smartconn.cc/tosone/ra-plus/store/tables"
)

// State 状态存储
type State struct {
	vi                *voiceInteract.VoiceInteract
	insertID          string
	bookID            string
	uniqueID          string
	globalSetting     globalSetting
	questions         map[string]map[string][]questionSetting
	settingFinish     *sync.WaitGroup // 解析脚本期间不能有其他动作
	current           *Current
	tips              []tables.Tip
	continuousRight   []string
	continuousWrong   []string
	continuousSilence []string
}

// Current 当前本页的状态
type Current struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	pageID       string
	question     []questionSetting
	questionID   string
	award        uint
	keywords     []string
	sttResult    string
	retryTimes   uint
	responseTime float64
}

func New() (state *State, err error) {
	state = new(State)
	if state.vi, err = voiceInteract.NewVoiceInteract(false, ""); err != nil {
		return
	}
	return
}

// Destroy 销毁
func (s *State) Destroy() {
	s.vi.Exit()
}

// Stop 停止当前的脚本执行
func (s *State) Stop() (err error) {
	return
}

// Setting 将一些附加的参数设置进来
func (s *State) Setting(insert, bookID, uniqueID string) (err error) {
	s.settingFinish = new(sync.WaitGroup)
	s.settingFinish.Add(1)
	defer s.settingFinish.Done()

	s.insertID = insert
	s.bookID = bookID
	s.uniqueID = uniqueID

	s.current = new(Current)

	var script = new(tables.InteractScript)
	if script, err = new(tables.InteractScript).FindByBookID(bookID); err != nil {
		return
	}
	s.questions = map[string]map[string][]questionSetting{}
	s.globalSetting = globalSetting{}
	if s.globalSetting, s.questions, err = s.parse(script.Manifest); err != nil {
		return
	}
	if s.tips, err = new(tables.Tip).GetAll(); err != nil {
		return
	}
	return
}
