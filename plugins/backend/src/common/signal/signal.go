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
package signal

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	chanWaitExitSignal = make(chan bool, 1)
)

func SetSignal() {
	go func() {
		chanSignal := make(chan os.Signal, 1)
		signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		for {
			signal := <-chanSignal
			switch signal {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP:
				chanWaitExitSignal <- true
				break
			default:
			}
		}
	}()
}

func WaitForSignal(_i32Timeout int) int {
	select {
	case bExit := <-chanWaitExitSignal:
		if bExit {
			return -1
		}
	case <-time.After(time.Second * time.Duration(_i32Timeout)):
	}
	return 0
}

//func WaitForEvent
