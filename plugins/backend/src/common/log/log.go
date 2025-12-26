package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	utils "cia/common/utils"
)

var lockMutex sync.Mutex

const (
	LOG_TYPE_ROTATE = 1
	LOG_TYPE_STDOUT = 2
	LOG_TYPE_DEBUG  = 4
)

type TLogObject struct {
	strName  string
	strPath  string
	iLogType int
	i32Level int32

	strModulePath  string
	iModulePathLen int

	i64Created      int64
	i32RotationCnt  int32
	i64RotationFreq int64

	i32CallStack int
}

func (obj *TLogObject) SetLog(_strName string, _strPath string, _iLogType int, _i32Level int32, _i32RotationCnt int32, _i64RotationFreq int64) {
	obj.strName = _strName
	obj.strPath = _strPath
	obj.i32Level = _i32Level
	obj.iLogType = _iLogType
	obj.i32RotationCnt = _i32RotationCnt
	obj.i64RotationFreq = _i64RotationFreq
	obj.strModulePath, _ = os.Getwd()
	obj.iModulePathLen = len(obj.strModulePath)

	obj.i32CallStack = 2

	utils.Mkdir(_strPath)
}

func (obj *TLogObject) SetCallStack(_iCall int) {
	obj.i32CallStack = _iCall
}

// WriteLog : write log text/json
func (obj *TLogObject) WriteF(_i32Level int32, _args ...interface{}) {

	var i64TmCurrent int64

	if obj.i32Level < _i32Level {
		return
	}

	var strTemp string
	for i := range _args {
		strTemp += fmt.Sprintf("%v", _args[i])
	}
	strText := obj._compileLogstring(strTemp[1 : len(strTemp)-1])

	lockMutex.Lock()
	defer lockMutex.Unlock()

	i64TmCurrent = time.Now().Unix()
	if obj._tryRotate(i64TmCurrent) == false {
	}

	strFilePath := fmt.Sprintf("%v/%v.txt", obj.strPath, obj.strName)
	if (obj.iLogType & LOG_TYPE_ROTATE) == LOG_TYPE_ROTATE {
		_writeFile(strFilePath, strText)
	} else if (obj.iLogType & LOG_TYPE_STDOUT) == LOG_TYPE_STDOUT {
		_writeConsole(strText)
	}
}

// WriteLog : write log text/json
func (obj *TLogObject) Write(_i32Level int32, _strLog string) {

	var i64TmCurrent int64

	if obj.i32Level < _i32Level {
		return
	}

	strText := obj._compileLogstring(_strLog)

	lockMutex.Lock()
	defer lockMutex.Unlock()

	i64TmCurrent = time.Now().Unix()
	if obj._tryRotate(i64TmCurrent) == false {
	}

	strFilePath := fmt.Sprintf("%v/%v.txt", obj.strPath, obj.strName)
	if (obj.iLogType & LOG_TYPE_ROTATE) == LOG_TYPE_ROTATE {
		_writeFile(strFilePath, strText)
	} else if (obj.iLogType & LOG_TYPE_STDOUT) == LOG_TYPE_STDOUT {
		_writeConsole(strText)
	}
}

func (obj *TLogObject) GetLogName() string {
	return obj.strName
}
