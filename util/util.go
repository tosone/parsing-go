package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"gopkg.in/satori/go.uuid.v1"
)

func times(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}

// PaddingLeft left-pads the string with pad up to len runes
// len may be exceeded if
func PaddingLeft(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

// PaddingRight right-pads the string with pad up to len runes
func PaddingRight(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

// UUIDWithDate 获取一个含有日期的 uuid
func UUIDWithDate() string {
	return UUID() + "_" + Now()
}

// UUIDPure 获取一个标准 uuid
func UUID() string {
	return strings.Join(strings.Split(uuid.NewV4().String(), "-"), "")
}

// Now 当前时间
func Now() string {
	return time.Now().Format("20060102150405")
}

// EnsureDir 确保文件夹存在
func EnsureDir(d string) error {
	if !com.IsDir(d) {
		return os.MkdirAll(d, 0755)
	}
	return nil
}

// ParentDir 获取父级目录
func ParentDir(p string) (s string) {
	if len(p) == 0 {
		s = "/"
	} else if len(p) == 1 && p == "/" {
		s = "/"
	} else if p[:1] != "/" {
		s = "/"
	} else if p[len(p)-1:] == "/" {
		list := strings.Split(p[:len(p)-1], "/")
		s = strings.Join(list[:len(list)-1], "/")
		if s == "" {
			s = "/"
		}
	} else {
		s = filepath.Dir(p)
	}
	return
}

// RandSelect 随机选择
func RandSelect(strList []string) string {
	if len(strList) == 0 {
		return ""
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return strList[rand.Intn(len(strList))]
}

// CheckCtx 检查 ctx 是否存活
func CheckCtx(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	default:
	}
	return true
}

// MD5Verify 用 MD5 验证一个文件
func MD5Verify(file, hash string) (verify bool, err error) {
	if !com.IsFile(file) {
		err = fmt.Errorf("No such a file: %s", file)
		return
	}
	var fileObj *os.File
	if fileObj, err = os.Open(file); err != nil {
		return
	}
	var data = make([]byte, 10240)
	var n int
	h := md5.New()
	for {
		n, err = fileObj.Read(data)
		if n != 0 {
			h.Write(data[:n])
		}
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
	}
	verify = hex.EncodeToString(h.Sum(nil)) == hash
	return
}

// HashString 哈希一个字符串
func HashString(str string) (hash []byte, err error) {
	h := md5.New()
	if _, err = h.Write([]byte(str)); err != nil {
		return
	}
	hash = h.Sum(nil)
	return
}

func HashFile(file string) (hash string, err error) {
	if !com.IsFile(file) {
		err = fmt.Errorf("No such a file: %s", file)
		return
	}
	var fileObj *os.File
	if fileObj, err = os.Open(file); err != nil {
		return
	}
	var data = make([]byte, 10240)
	var n int
	h := md5.New()
	for {
		n, err = fileObj.Read(data)
		if n != 0 {
			h.Write(data[:n])
		}
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
	}
	hash = hex.EncodeToString(h.Sum(nil))
	return
}

func Has(list []string, key string) bool {
	for _, l := range list {
		if l == key {
			return true
		}
	}
	return false
}

// CRC64ECMA CRC 校验
func CRC64ECMA(file string) (checksum string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err != nil {
		return
	}
	var s uint64
	s = crc64.Checksum(b, crc64.MakeTable(crc64.ECMA))
	checksum = strconv.FormatUint(s, 10)
	return
}
