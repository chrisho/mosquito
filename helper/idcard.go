package helper

import (
	"time"
	"github.com/chrisho/mosquito/validator"
	"strconv"
	"math"
)

func CalcIdcardAge(idcard string) int {
	if ! validator.IsIdcard(idcard) {
		return 0
	}

	idcardDateStr := string([]byte(idcard)[6:14])
	nowTimeDateStr := time.Now().Format("20060102")
	idcardDateInt, _ := strconv.Atoi(idcardDateStr)
	nowTimeDateInt, _ := strconv.Atoi(nowTimeDateStr)

	return int(math.Floor(float64(nowTimeDateInt-idcardDateInt)) / 10000)
}
