package helper

import "time"

const (
	YYYYMMDDHHIISS = "2006-01-02 15:04:05"
	YYYYMMDDHHII   = "2006-01-02 15:04"
	YYYYMMDDHH     = "2006-01-02 15"
	YYYYMMDD       = "2006-01-02"
)

func TimeNow() time.Time {
	return time.Now()
}

// get time now unixtime == php time()
func UnixTime() int32 {
	return int32(time.Now().Unix())
}

// get time now date format YYYY-MM-DD HH:II:SS
func DateTime() string {
	return time.Now().Format(YYYYMMDDHHIISS)
}

// RFC3339 time string to format Y-m-d H:i,Y-m-d H ...
func RFC3339Time2Date(timeStr, dateLayout string) string {
	timestamp, _ := time.ParseInLocation(time.RFC3339, timeStr, time.Local)

	return timestamp.Format(dateLayout)
}

// change date to format Y-m-d H:i,Y-m-d H ...
func Date2Date(inputDateLayout, inputDate, resultDateLayout string) string {
	timestamp, _ := time.ParseInLocation(inputDateLayout, inputDate, time.Local)

	return timestamp.Format(resultDateLayout)
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
	return t2.UnixNano() / 1e6 - t1.UnixNano() / 1e6
}

func TimeSince(start time.Time) time.Duration {
	return time.Since(start)
}
