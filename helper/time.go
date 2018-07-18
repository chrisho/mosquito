package helper

import (
	"time"
)

//const timezone = "Asia/ShangHai"
//
//func init() {
//	// 设置时间区
//	if loc, err := time.LoadLocation(timezone); err != nil {
//		println(err)
//	} else {
//		time.Local = loc
//	}
//}

const (
	YYYYMMDDHHIISS = "2006-01-02 15:04:05"
	YYYYMMDDHHII   = "2006-01-02 15:04"
	YYYYMMDDHH     = "2006-01-02 15"
	YYYYMMDD       = "2006-01-02"
	RFC3339        = time.RFC3339
)

func TimeNow() time.Time {
	return time.Now()
}

func TodayDate() string {
	return time.Now().Format(YYYYMMDD)
}

func TodayUnix() int64 {
	dateTime, _ := time.ParseInLocation(YYYYMMDD, TodayDate(), time.Local)

	return dateTime.Unix()
}

// get time now unixtime == php time()
func UnixTime() int32 {
	return int32(time.Now().Unix())
}

// get time now date format YYYY-MM-DD HH:II:SS
func DateTime() string {
	return time.Now().Format(YYYYMMDDHHIISS)
}

// 这月的某天
func ThisMonthDayTime(timeTime time.Time, whichDay int) time.Time {
	if whichDay <= 0 {
		whichDay = 1
	}
	return time.Date(timeTime.Year(), timeTime.Month(), whichDay, 0, 0, 0, 0, time.Local)
}

// 下月1号时间
func NextMonth01Time(timeTime time.Time) time.Time {
	return time.Date(timeTime.Year(), timeTime.Month(), 1, 0, 0, 0, 0, time.Local).
		AddDate(0, 1, 0)
}

// RFC3339 time string to format Y-m-d H:i,Y-m-d H ...
func RFC3339Time2Date(timeStr, dateLayout string) string {
	dateTime, _ := time.ParseInLocation(time.RFC3339, timeStr, time.Local)

	return dateTime.Format(dateLayout)
}

// change date to format Y-m-d H:i,Y-m-d H ...
func Date2Date(inputDateLayout, inputDate, resultDateLayout string) string {
	dateTime, _ := time.ParseInLocation(inputDateLayout, inputDate, time.Local)

	return dateTime.Format(resultDateLayout)
}

// date to time.Time
func Date2Time(date, dateFormat string) time.Time {
	timeTime, _ := time.ParseInLocation(dateFormat, date, time.Local)

	return timeTime
}

// Timestamp to time.Time
func Timestamp2Time(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

//  date to timestamp
func Date2Timestamp(dateLayout, date string) int64 {
	localTime, _ := time.ParseInLocation(dateLayout, date, time.Local)

	return localTime.Unix()
}

// timestamp to format Y-m-d H:i:s
func Timestamp2Date(timestamp int64, dateLayout string) string {
	return time.Unix(timestamp, 0).Format(dateLayout)
}

// 两个时间相差毫秒数
func TwoTimeMillisecond(t1, t2 time.Time) int64 {
	return t2.UnixNano()/1e6 - t1.UnixNano()/1e6
}

func TimeSince(start time.Time) time.Duration {
	return time.Since(start)
}

func DateSlice() []byte {
	var dateSlice []byte
	dateStr := DateTime()
	dateStrSlice := []byte(dateStr)

	for _, r := range dateStrSlice {
		if r >= '0' && r <= '9' {
			dateSlice = append(dateSlice, r)
		}
	}

	return dateSlice
}

func DateYYYYMMDDHHIISS() string {
	dateSlice := DateSlice()

	return string(dateSlice)
}

func DateYYMMDDHHIISS() string {
	dateSlice := DateSlice()

	return string(dateSlice[2:len(dateSlice)-1])
}
