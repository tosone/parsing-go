package interact

import "errors"

var errVoiceStartFail = errors.New("Voice start fail")         // 在规定的时间之前没有任何声音
var errVoiceDetectTimeout = errors.New("Voice detect timeout") // 声音识别已经达到最长时间

var errBackWordsVoid = errors.New("Said nothing") // 声音识别的内容是空
var errBackWordsWrong = errors.New("Said wrong")  // 说话识别的内容不在设置的列表中

var errVoiceDetectSuccess = errors.New("Voice detect succ") // 语音成功识别

var errBadNetwork = errors.New("upload voice data but bad network break it") // 弱网络导致上传声音数据失败
