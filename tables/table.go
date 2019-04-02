package tables

// FollowScript 跟读脚本
type FollowScript struct {
	BookID   string // 一本书的唯一 ID
	Manifest []byte // 脚本内容
	AudioRef []byte // 标准音频的 ref
}

// Tip 提示音频表
type Tip struct {
	Hash string // 音频的 hash 值
	Path string // 音频的路径
}

// FollowBaseScript 跟读的基础脚本
type FollowBaseScript struct {
	ScriptID string // 脚本编号
	Version  string // 脚本版本
	Manifest []byte // 脚本内容
}
