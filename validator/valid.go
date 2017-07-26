package validator

import (
	"github.com/asaskevich/govalidator"
	"log"
	"strconv"
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
			log.Println("phone string to int64 err : ", err)
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
