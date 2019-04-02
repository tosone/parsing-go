package tips

import (
	"fmt"
	"path"
	"path/filepath"
)

var tipPrefix = ""

// Get 获取提示音频
func Get(name string) (tip string) {
	tip = path.Join(tipPrefix, name+".ogg")
	return
}

// Gets 获取提示音频列表
func Gets(name string) (tips []string, err error) {
	tips, err = filepath.Glob(fmt.Sprintf("%s/%s*.ogg", tipPrefix, name))
	return
}
