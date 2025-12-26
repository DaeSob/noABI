// V2.0.0 By XeN
package preference

import (
	"time"
)

func onSetTimeZone() {
	inst := GetInstance()

	t1 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.Local)
	t2 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	inst.timeZone = t2.Unix() - t1.Unix()
}

func GetTimeZone() int64 {
	inst := GetInstance()

	return inst.timeZone
}
