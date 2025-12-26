package utils

import (
	"time"
)

func Now() time.Time {
	now := time.Now()
	secs := now.Unix()
	return time.Unix(secs, 0)
}

func CutMillisecond(_t time.Time) time.Time {
	secs := _t.Unix()
	return time.Unix(secs, 0)
}

func LocalTimeToGMT(_t time.Time) time.Time {
	gmtLoc, _ := time.LoadLocation("GMT")
	return _t.In(gmtLoc)
}

func LocalTimeToGMTString(_t time.Time) string {
	gmtLoc, _ := time.LoadLocation("GMT")
	gmt := _t.In(gmtLoc)
	return gmt.Format("Mon, _2 Jan 2006 15:04:05 GMT")
}

func GMTStringToLocalTime(_gmt string) time.Time {
	if len(_gmt) < 25 {
		t, _ := time.Parse("Mon, _2 Jan 2006 15:04:05", "")
		return t
	}
	t, _ := time.Parse("Mon, _2 Jan 2006 15:04:05", _gmt[:25])
	loc, _ := time.LoadLocation("Local")
	return t.In(loc)
}
