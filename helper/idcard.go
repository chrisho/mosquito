package helper

import "time"

func CalcIdcardAge(idcard string) int {
	dateByte := []byte(TrimStringSpace(idcard))[6:14]

	dateStr := string(dateByte[0:4]) + "-" + string(dateByte[4:6]) + "-" + string(dateByte[6:8])

	timeTime := Date2Time(dateStr, YYYYMMDD)
	timeNow := time.Now()

	age := timeNow.Year() - timeTime.Year()

	if timeNow.Month() - timeTime.Month() < 0 {
		age -= 1
	}

	if timeNow.Month() - timeTime.Month() == 0 && timeNow.Day() - timeTime.Day() < 0 {
		age -= 1
	}
	return age
}
