package model

// Tip 提示音频表
type Tip struct {
	Hash string `json:"hash"` // 音频的 hash 值
	Path string `json:"path"` // 音频的路径
}
