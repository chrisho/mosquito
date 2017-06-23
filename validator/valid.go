package validator

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.CustomTypeTagMap.Set("mobilePhone", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {

		phone := i.(int64)

		return phone > 10000000000 && phone < 19999999999
	}))

}