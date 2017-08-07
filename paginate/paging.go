package paginate

import (
	"strings"
	"github.com/chrisho/mosquito/utils"
	"reflect"
)

const (
	PagingSize      = 10
	PagingByPrimary = iota
	PagingByNumber
)

// 获取分页选项
func GetPagingOptions(in *PageOptions, PagingMode int) (offset, limit int32) {

	SetPagingDefaultOptions(in)

	switch PagingMode {

	case PagingByPrimary:
		offset, limit = GetPagingModeByPrimaryOptions(in)
	case PagingByNumber:
		offset, limit = GetPagingModeByNumberOptions(in)
	}

	return
}

// 设置 : 默认每页 10 条，页码 第 1 页
func SetPagingDefaultOptions(in *PageOptions) *PageOptions {

	// set default pageSize ：
	if in.PageSize < 1 {
		in.PageSize = int32(PagingSize)
	}

	// set default first page : 1
	if in.PageNumber < 1 {
		in.PageNumber = 1
	}
	// 设置默认查询字段、排序
	in.SortField, in.SortFieldTo = SetPagingModeByPrimarySelectFieldAndSort(in.SortField, in.SortFieldTo)

	return in
}

// 默认排序
func SetPagingModeByPrimarySelectFieldAndSort(SortField, SortFieldTo string) (field string, sort string) {
	sort = "desc"

	if field = strings.Trim(SortField, " "); field == "" {
		field = "id"
	}

	field = utils.SnakeString(field)

	SortFieldTo = strings.ToLower(strings.Trim(SortFieldTo, " "))

	if SortFieldTo == "asc" {
		sort = SortFieldTo
	}

	return field, sort
}

// structPointer 必须是 struct 的 指针
func PagingOptionsFieldNameIsValid(structPointer interface{}, fieldName string) bool {

	sElem := reflect.ValueOf(structPointer).Elem()

	return sElem.FieldByName(utils.CamelString(fieldName)).IsValid()
}

// 页码分页模式选项
func GetPagingModeByNumberOptions(in *PageOptions) (offset, limit int32) {

	offset = in.PageSize * (in.PageNumber - 1)
	limit = in.PageSize

	return
}

/*
 * select * from users where id > ? order by id asc limit 0,PageSize;
 * select * from users where id < ? order by id desc limit 0,PageSize;
 * if _, ok := requestParams["SortValue"]; ok {
 *		sortField := SellPointLimitTable + "." + requestParams["SortField"].(string)
 *		sortFieldTo := requestParams["SortFieldTo"].(string)
 *		orderBy = sortField + " " + sortFieldTo
 *		switch sortFieldTo {
 *		case "asc":
 *			sql = sql.Where(sortField + " > ?", requestParams["SortValue"])
 *		default:
 *			if requestParams["SortValue"].(int64) > 0 {
 *				sql = sql.Where(sortField + " < ?", requestParams["SortValue"])
 *			}
 *		}
 *	}
 */
// 主键分页模式选项
func GetPagingModeByPrimaryOptions(in *PageOptions) (offset, limit int32) {

	offset = 0
	limit = in.PageSize

	return
}

// Set Paging Result
func SetPagingResult(in *PageOptions, TotalRecords int32, SortValue int64) (paginate PageResult) {

	paginate.TotalRecords = TotalRecords

	if paginate.TotalRecords%in.PageSize == 0 {
		paginate.TotalPages = paginate.TotalRecords / in.PageSize
	} else {
		paginate.TotalPages = paginate.TotalRecords/in.PageSize + 1
	}

	paginate.PageSize = in.PageSize
	paginate.PageNumber = in.PageNumber
	paginate.SortValue = SortValue

	return
}
