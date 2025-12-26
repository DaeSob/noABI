package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// write file from string
func WriteFileFromString(_strPath string, _strText string) {
	f, e := os.OpenFile(_strPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	f.WriteString(_strText)
}

func GetPath() string {
	path, _ := os.Getwd()
	return path
}

func ExistPath(_strPath string) bool {
	if _, err := os.Stat(_strPath); err != nil {
		return false
	}
	return true
}

func GetFileSize(_filePath string) int64 {
	fileInfo, err := os.Stat(_filePath)
	if err != nil {
		panic(err)
	}

	return fileInfo.Size()
}

func Mkdir(_strPath string) error {
	if _, serr := os.Stat(_strPath); serr != nil {
		merr := os.MkdirAll(_strPath, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func Remove(_strPath string) error {
	return os.Remove(_strPath)
}

func RemoveAll(_strPath string) error {
	return os.RemoveAll(_strPath)
}

func Rename(_strFrom string, _strTo string) error {
	return os.Rename(_strFrom, _strTo)
}

func ChangeFileTime(_strPath string, _tmAccess time.Time, _tmMod time.Time) error {
	return os.Chtimes(_strPath, _tmAccess, _tmMod)
}

// ///////////////////////////////////////////////////////////////////////
// SFolder_r : get file list include sub-dir
// ///////////////////////////////////////////////////////////////////////
func SFolder_r(_fileList *[]string, _path string, _whiteList []string) error {
	files, err := ioutil.ReadDir(_path)
	if nil != err {
		return err
	}

	// get file list
	for _, f := range files {
		// when file attr is directory,
		if true == f.IsDir() {
			// recursively running
			SFolder_r(_fileList, fmt.Sprint(_path, "/", f.Name()), _whiteList)
			continue
		}

		// check white list extension
		absFile := fmt.Sprint(_path, "/", f.Name())
		if true == extChecker(absFile, _whiteList) {
			*_fileList = append(*_fileList, absFile)
		}
	}

	return nil
}

// get_current_files : get file list except sub-dir
func sfolder(_fileList *[]string, _path string) error {
	files, err := ioutil.ReadDir(_path)
	if nil != err {
		return err
	}

	for _, f := range files {
		if true == f.IsDir() {
			continue
		}
		*_fileList = append(*_fileList, fmt.Sprint(_path, "/", f.Name()))
	}
	return nil
}

// extChecker : file extention checker
func extChecker(_path string, _whiteList []string) bool {
	for _, v := range _whiteList {
		if v == filepath.Ext(strings.TrimSpace(_path)) {
			return true
		}
	}
	return false
}

type TFileInfo struct {
	Path    string
	Name    string
	Size    int64
	IsDir   bool
	ModTime time.Time
}

// get_current_files : get file list except sub-dir
func SFolder(_strPath string, _sliceFileInfo *[]TFileInfo) error {

	files, err := ioutil.ReadDir(_strPath)
	if nil != err {
		return err
	}

	for _, f := range files {

		var fileInfo TFileInfo

		fileInfo.Path = _strPath
		fileInfo.Name = f.Name()
		fileInfo.Size = f.Size()
		fileInfo.IsDir = f.IsDir()
		fileInfo.ModTime = f.ModTime()

		*_sliceFileInfo = append(*_sliceFileInfo, fileInfo)
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////////
// FileCopy_s : file copy with temporary extension
// ///////////////////////////////////////////////////////////////////////
func FileCopy_s(_src string, _dst string) (int64, error) {
	sourceFileStat, err := os.Stat(_src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", _src)
	}

	source, err := os.Open(_src)
	if err != nil {
		return 0, err
	}

	// add temporary extention ".lock"
	tmp := _dst + ".lock"
	destination, err := os.Create(tmp)
	if err != nil {
		return 0, err
	}
	nBytes, err := io.Copy(destination, source)
	source.Close()
	destination.Close()
	// when finished copy, MoveFile : tmp to dst
	err = os.Rename(tmp, _dst)

	return nBytes, err
}
