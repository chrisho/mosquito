package validator

import (
	"github.com/asaskevich/govalidator"
	"strconv"
	"github.com/sirupsen/logrus"
	"regexp"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	mobilePhone()
}

func mobilePhone()  {
	govalidator.CustomTypeTagMap.Set("mobilePhone", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {

		phoneStr := i.(string)

		// 字符串转int64
		phone, err := strconv.ParseInt(phoneStr, 10, 64)

		if err != nil {
			logrus.Error("phone string to int64 err : ", err)
			return false
		}

		return phone > 10000000000 && phone < 19999999999
	}))
}

/*
Here is a list of available validators for struct fields (validator - used function):

"email":          IsEmail,
"url":            IsURL,
"dialstring":     IsDialString,
"requrl":         IsRequestURL,
"requri":         IsRequestURI,
"alpha":          IsAlpha,
"utfletter":      IsUTFLetter,
"alphanum":       IsAlphanumeric,
"utfletternum":   IsUTFLetterNumeric,
"numeric":        IsNumeric,
"utfnumeric":     IsUTFNumeric,
"utfdigit":       IsUTFDigit,
"hexadecimal":    IsHexadecimal,
"hexcolor":       IsHexcolor,
"rgbcolor":       IsRGBcolor,
"lowercase":      IsLowerCase,
"uppercase":      IsUpperCase,
"int":            IsInt,
"float":          IsFloat,
"null":           IsNull,
"uuid":           IsUUID,
"uuidv3":         IsUUIDv3,
"uuidv4":         IsUUIDv4,
"uuidv5":         IsUUIDv5,
"creditcard":     IsCreditCard,
"isbn10":         IsISBN10,
"isbn13":         IsISBN13,
"json":           IsJSON,
"multibyte":      IsMultibyte,
"ascii":          IsASCII,
"printableascii": IsPrintableASCII,
"fullwidth":      IsFullWidth,
"halfwidth":      IsHalfWidth,
"variablewidth":  IsVariableWidth,
"base64":         IsBase64,
"datauri":        IsDataURI,
"ip":             IsIP,
"port":           IsPort,
"ipv4":           IsIPv4,
"ipv6":           IsIPv6,
"dns":            IsDNSName,
"host":           IsHost,
"mac":            IsMAC,
"latitude":       IsLatitude,
"longitude":      IsLongitude,
"ssn":            IsSSN,
"semver":         IsSemver,
"rfc3339":        IsRFC3339,
"ISO3166Alpha2":  IsISO3166Alpha2,
"ISO3166Alpha3":  IsISO3166Alpha3,

Validators with parameters

"range(min|max)": Range,
"length(min|max)": ByteLength,
"runelength(min|max)": RuneLength,
"matches(pattern)": StringMatches,
"in(string1|string2|...|stringN)": IsIn,
*/

// IsPhoneNumber
func IsPhoneNumber(value string) bool {
	reg := `^1([3-9][0-9]|14[57]|5[^4])\d{8}$`

	return regexp.MustCompile(reg).MatchString(value)
}

// IsFixedPhone
func IsFixedPhone(value string) bool {
	reg := `^(0[0-9]{2,3}-)?([2-9][0-9]{6,7})+(-[0-9]{1,4})?$`

	return regexp.MustCompile(reg).MatchString(value)
}

// bank account
func IsBankAccount(value string) bool {
	reg := `^([1-9]{1})(\d{15}|\d{18})$`

	return regexp.MustCompile(reg).MatchString(value)
}

// idcard
func IsIdcard(value string) bool {
	if len(value) == 15 {
		reg := `^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}[0-9Xx]$`
		return regexp.MustCompile(reg).MatchString(value)
	} else {
		reg := `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
		return regexp.MustCompile(reg).MatchString(value)
	}
	return false
}

// url 
func IsUrl(value string) bool {
	return govalidator.IsURL(value)
}

// chinese
func IsChinese(value string) bool {
	reg := `^[\p{Han}]+$`

	return regexp.MustCompile(reg).MatchString(value)
}

// chinese and letter and number
func IsChineseAndLetterAndNumber(value string) bool {
	reg := `^[a-zA-Z0-9\p{Han}]+$`

	return regexp.MustCompile(reg).MatchString(value)
}

// postcode
func IsPostcode(value string) bool {
	reg := `^[1-9][0-9]{5}$`

	return regexp.MustCompile(reg).MatchString(value)
}

// IsNumeric
func IsNumeric(value string) bool {
	if value == "" {
		return false
	}
	return govalidator.IsNumeric(value)
}

// IsEmail
func IsEmail(value string) bool {
	return govalidator.IsEmail(value)
}

// IsImei
func IsImei(value string) bool {
	reg := `^[0-9A-Za-z]+$`

	return regexp.MustCompile(reg).MatchString(value)
}