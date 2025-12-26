package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	utils "cia/common/utils"
)

func (obj *TLogObject) _tryRotate(_i64TmCurrent int64) bool {

	if obj.i32RotationCnt == 0 {
		return false
	}

	strMetaFilePath := fmt.Sprintf("%v/%v.meta", obj.strPath, obj.strName)

	if obj.i64Created == 0 {
		if utils.ExistPath(strMetaFilePath) == true {
			strCTime, _ := _readFromMetaFile(strMetaFilePath)
			obj.i64Created, _ = utils.StringToI64(strCTime)
		} else {
			obj.i64Created = _i64TmCurrent
			_writeToMetaFile(utils.I64ToString(obj.i64Created), strMetaFilePath)
		}
	}

	if obj.i64Created > 0 {
		i64TempCTime := (obj.i64Created / obj.i64RotationFreq) * obj.i64RotationFreq
		i64TempCurTime := (_i64TmCurrent / obj.i64RotationFreq) * obj.i64RotationFreq

		if i64TempCTime < i64TempCurTime {
			//Try Rotate
			strLastFilePath := fmt.Sprintf("%v/%v.txt.%v", obj.strPath, obj.strName, obj.i32RotationCnt)

			if utils.ExistPath(strLastFilePath) == true {
				utils.Remove(strLastFilePath)
			}

			var i32 int32
			for i32 = obj.i32RotationCnt - 1; i32 >= 0; i32-- {
				var strNextFilePath string
				var strPreFilePath string
				if i32 == 0 {
					strPreFilePath = fmt.Sprintf("%v/%v.txt", obj.strPath, obj.strName)
					strNextFilePath = fmt.Sprintf("%v/%v.txt.%v", obj.strPath, obj.strName, i32+1)
					//Metafile 갱신
					obj.i64Created = _i64TmCurrent
					_writeToMetaFile(fmt.Sprintf("%v", obj.i64Created), strMetaFilePath)
				} else {
					strPreFilePath = fmt.Sprintf("%v/%v.txt.%v", obj.strPath, obj.strName, i32)
					strNextFilePath = fmt.Sprintf("%v/%v.txt.%v", obj.strPath, obj.strName, i32+1)
				}

				if utils.ExistPath(strPreFilePath) == true {
					err := utils.Rename(strPreFilePath, strNextFilePath)
					if err != nil {
						//TODO
					}
				}
			}
			return true
		}
	} else {
		//Crate New Log file
		fmt.Println("New LogFile.....")
	}
	return false
}

// CompileLogstring : Compile Log string from args
func (obj *TLogObject) _compileLogstring(_strLog string) string {
	var strText string
	localTime := time.Now().Local()

	strText += localTime.Format("2006-01-02 15:04:05") + "\t"
	if (obj.iLogType & LOG_TYPE_DEBUG) == 4 {
		_, strFile, iLine, _ := runtime.Caller(obj.i32CallStack)
		strText += fmt.Sprintf("/%s %d\t", filepath.Base(strFile), iLine)
	}
	return strText + _strLog + "\r\n"
}

// WriteToMetaFile : write to meta file
func _writeToMetaFile(_strText string, _strLogFile string) {
	f, e := os.OpenFile(_strLogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	f.WriteString(_strText)
}

// ReadFromMetaFile : read from meta file
func _readFromMetaFile(_strLogMetaFile string) (string, error) {
	f, e := os.OpenFile(_strLogMetaFile, os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		fmt.Println(e)
		return "", e
	}
	defer f.Close()
	buffer := make([]byte, 10)
	_, e = f.Read(buffer)
	if e != nil {
		return "", e
	}
	return string(buffer), nil
}

// WriteFile : write log at file
func _writeFile(_strLogFile string, _strText string) {
	utils.WriteFileFromString(_strLogFile, _strText)
}

// WriteConsole : write log at stanadard output buffer
func _writeConsole(_text string) {
	fmt.Println(_text)
}
