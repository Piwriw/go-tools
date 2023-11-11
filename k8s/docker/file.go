package docker

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

var (
	CfgFileMD5 string
)

// GenMd5 为文件生成 MD5 值
func GenMd5(absPath string) (string, error) {
	path, err := filepath.Abs(absPath)
	if err != nil {
		return "", errors.Wrap(err, "GenMd5 出错：%v")
	}

	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "GenMd5 出错：%v")
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "copy bitwise failed")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// DirFileIsExist 判断文件或者目录是否存在
func DirFileIsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// MakeDirAll 新建一个文件夹，首先会判断目录是否存在，如果不存在，那么就会创建该目录，否则什么都不做
func MakeDirAll(path string) error {
	tf := DirFileIsExist(path)
	if !tf {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSubDirNameList 返回一个目录中子目录都路径列表，（判断软链接）
func GetSubDirNameList(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	subDirList := make([]string, 0)
	for _, f := range files {
		if f.IsDir() || !f.Mode().Type().IsRegular() {
			subDir := filepath.Join(dirPath, f.Name())
			subDirList = append(subDirList, subDir)
		}
	}
	return subDirList, nil
}

// FileNumber 返回给定文件夹中文件的数量
func FileNumber(path string) (int, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	fileList, _ := ioutil.ReadDir(path + "/" + files[0].Name())
	sum := 0
	for _, f := range fileList {
		if !f.IsDir() {
			sum++
		}
	}
	return sum, nil
}

// FileDetectedInfo 返回给定文件夹中文件的识别数量，例如： map[car:1 people:1 people_car:1]
func FileDetectedInfo(dirPath string) (map[string]int, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	// fmt.Println(files[0])

	fileList, _ := ioutil.ReadDir(dirPath + "/" + files[0].Name())
	data := map[string]int{}

	var fileNameList []string
	for _, f := range fileList {
		dirName := f.Name()
		// 获取去后缀的文件名
		fileprefix := path.Base(dirName)[0 : len(path.Base(dirName))-len(path.Ext(dirName))]
		// 去除文件名数字部分
		i := len(fileprefix) - 1
		for ; i > 0; i-- {
			if fileprefix[i] >= '0' && fileprefix[i] <= '9' {
				continue
			}
			break
		}
		fileNameList = append(fileNameList, fileprefix[0:i+1])
	}

	// 检测文件名、数量
	for i := range fileNameList {
		data[fileNameList[i]]++
	}

	return data, nil
}

// MapToJson 将map[string]int数据类型转换成string数据类型
func MapToJson(param map[string]int) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// 转换时间
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// 将数组 filename 按指定大小进行分隔
func GroupSplit(filename []string, subGroupLength int64) [][]string {
	max := int64(len(filename))
	var segmens = make([][]string, 0)
	quantity := max / subGroupLength
	remainder := max % subGroupLength
	i := int64(0)
	for i = int64(0); i < quantity; i++ {
		segmens = append(segmens, filename[i*subGroupLength:(i+1)*subGroupLength])
	}
	if quantity == 0 || remainder != 0 {
		segmens = append(segmens, filename[i*subGroupLength:i*subGroupLength+remainder])
	}
	return segmens
}
