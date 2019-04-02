package follow

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"smartconn.cc/tosone/parsing-go/audio"
	"smartconn.cc/tosone/parsing-go/model"
	"smartconn.cc/tosone/parsing-go/tables"
)

var baseScriptID int // 上次选择的 base 脚本的序号，以实现每次放上书所用的 base 脚本都不一样

var status bool

// Status 整个跟读流程中的对象结构体
type State struct {
	insertID       string               // 插入动作的 ID
	bookID         string               // 书的 ID
	uniqueID       string               // 书的唯一版本的 ID
	audioRefs      map[string]string    // 标准音频的 ref
	tips           []model.Tip          // 所有的提示音频
	baseScript     GlobalSetting        // 当前的 base 脚本
	script         *tables.FollowScript // 数据库中脚本的详细信息
	baseScriptRaw  string               // base 原始脚本
	setting        GlobalSetting        // 脚本的全局的配置
	allPages       WholePages           // 每一页的配置
	mixedSetting   GlobalSetting        // 经过混淆过后的全局配置
	nextPageTimer  *time.Timer          // 提醒下一页
	contentEntered bool                 // 是否进入过正文部分
	goal           bool                 // 是否已经达到目标
	extraGoal      bool                 // 是否已经达到任务之外的任务目标
	flower         int                  // 当前已经获得小红花的数量
	totalFlower    int                  // 总共需要的小红花数量
	task           int                  // 当前已经完成的任务数量
	extraFlower    int                  // 在任务之外的获得的小红花数量
	extraTask      int                  // 在任务之外已经完成的任务数量
	sentenceMap    map[string]int       // 已经跟读对的句子和次数的对印关系
	validFollow    uint                 // 有效跟读次数
	//eva            *asr.ProEva          // 对比音频相似程度的对象
	settingFinish *sync.WaitGroup // setting 设置完成
	ctxWg         *sync.WaitGroup // 当 ctx 被结束掉之后需要等到流程完成之后才能开始新的 ctx
	current       Current         // 当前页的状态
	audio         audio.Audio
}

type Current struct {
	ctx       context.Context    // 运行过程中的上下文
	ctxCancel context.CancelFunc // 运行过程中的上下文取消
	pageID    string             // 当前页页码
}

func New() (state *State, err error) {
	state = new(State)
	//state.eva = asr.NewProEva()
	return
}

// Destroy 销毁
func (s *State) Destroy() {
	//s.eva = nil
}

// Stop 停止当前的脚本执行
func (s *State) Stop() (err error) {
	return
}

func (s *State) Setting(insert, bookID, uniqueID string, script, baseScript string, audio audio.Audio, tips string) (err error) {
	s.settingFinish = new(sync.WaitGroup)
	s.settingFinish.Add(1)
	defer s.settingFinish.Done()

	s.insertID = insert
	s.bookID = bookID
	s.uniqueID = uniqueID

	if err = json.Unmarshal([]byte(tips), &s.tips); err != nil {
		return
	}

	//if s.script, err = new(tables.FollowScript).FindByBookID(s.bookID); err != nil {
	//	return
	//}
	s.script = &tables.FollowScript{BookID: bookID, Manifest: []byte(script)}
	s.baseScriptRaw = baseScript
	//if err = json.Unmarshal(s.script.AudioRef, &s.audioRefs); err != nil {
	//	return
	//}

	// 将脚本解析成全局配置和每页的配置
	var global = GlobalSetting{}
	if global, s.allPages, err = Parse(s.script.Manifest); err != nil {
		return
	}

	// 随机出来一个 base 脚本
	if s.baseScript, err = s.randomBaseScript(); err != nil {
		return
	}

	// 将本书的 global setting 和 base 脚本 mix 到一块
	if s.mixedSetting, err = mixGlobal(s.baseScript, global); err != nil {
		return
	}
	if err = s.calcTotalFlower(); err != nil {
		return
	}
	s.audio = audio
	return
}
