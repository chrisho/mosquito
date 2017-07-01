package validator

import (
	"github.com/asaskevich/govalidator"
	"strconv"
	"log"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.CustomTypeTagMap.Set("mobilePhone", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {

		phoneStr := i.(string)

		// 字符串转int64
		phone, err := strconv.ParseInt(phoneStr, 10, 64)

		if err != nil {
			log.Println("phone string to int64 err : ", err)
			return false
		}

		return phone > 10000000000 && phone < 19999999999
	}))

}