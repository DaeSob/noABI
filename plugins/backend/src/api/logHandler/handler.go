package logHandler

import (
	"fmt"
	"io/ioutil"
	"sync"

	log "cia/common/log"

	"gopkg.in/yaml.v2"
)

type TMyLog struct {
	LOGS struct {
		INITIALIZE struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"initialize"`
		ACCESS struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"access"`
		TRACE struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"trace"`
		DB struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"db"`
		QUEUEDDB struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"queuedDB"`
		FAILED struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"failed"`
		EXCEPTION struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"exception"`
		BLOCK struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"block"`
		PAYMENT struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"payment"`
		LOGCOLLECTOR struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"logCollector"`
		ENTRYPOINT struct {
			Path     string `yaml:"path"`
			Debug    bool   `yaml:"debug"`
			Level    int32  `yaml:"level"`
			RotaCnt  int32  `yaml:"rotacnt"`
			RotaFreq int64  `yaml:"rotafreq"`
		} `yaml:"entryPoint"`
	} `yaml:"logs"`
}

type TLogHandler struct {
	sliceLogObj []log.TLogObject
}

// using singleton
var instance *TLogHandler
var once sync.Once

func GetInstance() *TLogHandler {
	once.Do(func() {
		instance = &TLogHandler{}
	})
	return instance
}

func _getLog(_strLogName string) *log.TLogObject {
	inst := GetInstance()
	for i := range inst.sliceLogObj {
		if inst.sliceLogObj[i].GetLogName() == _strLogName {
			return &inst.sliceLogObj[i]
		}
	}
	return nil
}

func Initialize(
	_build string,
) error {
	inst := GetInstance()

	logPath := "./config/" + _build + "/log.yaml"
	byteData, err := ioutil.ReadFile(logPath)

	if err != nil {
		return err
	}

	myLog := TMyLog{}
	err = yaml.Unmarshal(byteData, &myLog)
	if err != nil {
		return err
	}

	//INITIALIZE LOG SETTING
	initLogObj := log.TLogObject{}
	var iLogType int
	if myLog.LOGS.INITIALIZE.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	initLogObj.SetLog("initialize", myLog.LOGS.INITIALIZE.Path, iLogType, myLog.LOGS.INITIALIZE.Level, myLog.LOGS.INITIALIZE.RotaCnt, myLog.LOGS.INITIALIZE.RotaFreq)
	initLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, initLogObj)

	//TRACE LOG SETTING
	traceLogObj := log.TLogObject{}
	if myLog.LOGS.TRACE.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	traceLogObj.SetLog("trace", myLog.LOGS.TRACE.Path, iLogType, myLog.LOGS.TRACE.Level, myLog.LOGS.TRACE.RotaCnt, myLog.LOGS.TRACE.RotaFreq)
	traceLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, traceLogObj)

	//DB LOG SETTING
	dbLogObj := log.TLogObject{}
	if myLog.LOGS.DB.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	dbLogObj.SetLog("db", myLog.LOGS.DB.Path, iLogType, myLog.LOGS.DB.Level, myLog.LOGS.DB.RotaCnt, myLog.LOGS.DB.RotaFreq)
	dbLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, dbLogObj)

	//TRACE LOG SETTING
	queuedDBLogObj := log.TLogObject{}
	if myLog.LOGS.QUEUEDDB.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	queuedDBLogObj.SetLog("queuedDB", myLog.LOGS.QUEUEDDB.Path, iLogType, myLog.LOGS.QUEUEDDB.Level, myLog.LOGS.QUEUEDDB.RotaCnt, myLog.LOGS.QUEUEDDB.RotaFreq)
	queuedDBLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, queuedDBLogObj)

	//FAILED LOG SETTING
	failedLogObj := log.TLogObject{}
	if myLog.LOGS.FAILED.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	failedLogObj.SetLog("failed", myLog.LOGS.FAILED.Path, iLogType, myLog.LOGS.FAILED.Level, myLog.LOGS.FAILED.RotaCnt, myLog.LOGS.FAILED.RotaFreq)
	failedLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, failedLogObj)

	//EXCEPTION LOG SETTING
	exceptionLogObj := log.TLogObject{}
	if myLog.LOGS.EXCEPTION.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	exceptionLogObj.SetLog("exception", myLog.LOGS.EXCEPTION.Path, iLogType, myLog.LOGS.EXCEPTION.Level, myLog.LOGS.EXCEPTION.RotaCnt, myLog.LOGS.EXCEPTION.RotaFreq)
	exceptionLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, exceptionLogObj)

	//ACCESS LOG SETTING
	accessLogObj := log.TLogObject{}
	if myLog.LOGS.ACCESS.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	accessLogObj.SetLog("access", myLog.LOGS.ACCESS.Path, iLogType, myLog.LOGS.ACCESS.Level, myLog.LOGS.ACCESS.RotaCnt, myLog.LOGS.ACCESS.RotaFreq)
	accessLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, accessLogObj)

	//Block LOG SETTING
	blockLogObj := log.TLogObject{}
	if myLog.LOGS.BLOCK.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	blockLogObj.SetLog("block", myLog.LOGS.BLOCK.Path, iLogType, myLog.LOGS.BLOCK.Level, myLog.LOGS.BLOCK.RotaCnt, myLog.LOGS.BLOCK.RotaFreq)
	blockLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, blockLogObj)

	//Payment LOG SETTING
	paymentLogObj := log.TLogObject{}
	if myLog.LOGS.PAYMENT.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	paymentLogObj.SetLog("payment", myLog.LOGS.PAYMENT.Path, iLogType, myLog.LOGS.PAYMENT.Level, myLog.LOGS.PAYMENT.RotaCnt, myLog.LOGS.PAYMENT.RotaFreq)
	paymentLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, paymentLogObj)

	//Log Collector LOG SETTING
	logCollectorLogObj := log.TLogObject{}
	if myLog.LOGS.LOGCOLLECTOR.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	logCollectorLogObj.SetLog("logCollector", myLog.LOGS.LOGCOLLECTOR.Path, iLogType, myLog.LOGS.LOGCOLLECTOR.Level, myLog.LOGS.LOGCOLLECTOR.RotaCnt, myLog.LOGS.LOGCOLLECTOR.RotaFreq)
	logCollectorLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, logCollectorLogObj)

	//Log Collector LOG SETTING
	logEntryPointLogObj := log.TLogObject{}
	if myLog.LOGS.ENTRYPOINT.Debug {
		iLogType = log.LOG_TYPE_ROTATE | log.LOG_TYPE_DEBUG
	} else {
		iLogType = log.LOG_TYPE_ROTATE
	}
	logEntryPointLogObj.SetLog("entryPoint", myLog.LOGS.ENTRYPOINT.Path, iLogType, myLog.LOGS.ENTRYPOINT.Level, myLog.LOGS.ENTRYPOINT.RotaCnt, myLog.LOGS.ENTRYPOINT.RotaFreq)
	logEntryPointLogObj.SetCallStack(3)
	inst.sliceLogObj = append(inst.sliceLogObj, logEntryPointLogObj)

	byteData = nil
	return err
}

func Write(_strLogName string, _i32Level int32, _args ...interface{}) bool {
	logObj := _getLog(_strLogName)

	if logObj == nil {
		return false
	}

	var strTemp string
	for i := range _args {
		if i > 0 {
			strTemp += fmt.Sprintf(" %v", _args[i])
		} else {
			strTemp += fmt.Sprintf("%v", _args[i])
		}
	}
	logObj.Write(_i32Level, strTemp)

	return true
}
