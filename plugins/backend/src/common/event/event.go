// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Hyper-Cluster Common
//
//
//
//																										2020.08.01
//																										DAESEOB.JEONG
//
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
package event

import (
	"time"
)

func CreateChannelEvent() chan bool {
	return make(chan bool, 1)
}

func SetChannelEvent(_bChannelEvent chan bool) {
	_bChannelEvent <- true
}

func ResetChannelEvent(_bChannelEvent chan bool) {
	_bChannelEvent <- false
}

func WaitFor2CHEvent(_bCH1 chan bool, _bCH2 chan bool, _i32MillisecondTimeout int) int {
	select {
	case bChEvent := <-_bCH1:
		if bChEvent == true {
			return 1
		} else {
			return -1
		}
	case bChEvent := <-_bCH2:
		if bChEvent == true {
			return 1
		} else {
			return -1
		}
	case <-time.After(time.Millisecond * time.Duration(_i32MillisecondTimeout)):
	}
	return 0
}
