package lib

import (
	"fmt"
	"time"
)

const (
	TimeSecond time.Duration = time.Second
	TimeMinute               = time.Minute
	TimeHour                 = time.Hour
	// TimeDay is not supported in default
	// time library. Therefore, create one
	TimeDay = 24 * time.Hour
)

// TimeHandler Offset is a number that will multiply TimUnit
// and the result will be added to the current time
type TimeHandler struct {
	//Offset   int
	TimeUnit time.Duration
}

// TimeStamp in 2018-07-19 09:53:22
// sring format
func (th TimeHandler) Now() string {
	return th.getTime(0, time.Second)
}

// make timestamp by adding any int and unit
// ex. 20, lib.TimeMinute
func (th TimeHandler) MakeTime(offset int, timeUnit time.Duration) string {
	return th.getTime(offset, timeUnit)
}

func (th TimeHandler) getTime(offset int, timeUnit time.Duration) string {
	//t := time.Now().Add(ts.TimeUnit * time.Duration(ts.Offset))
	t := time.Now().Add(timeUnit * time.Duration(offset))
	timeStamp := t.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("%v", timeStamp)
}
