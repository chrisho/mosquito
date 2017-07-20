package paginate

import (
	"strings"
)

var (
	DefaultPageSize int32 = 10 // 默认显示条数
	PaginateMode          = map[int32]string{0: "PageNumber", 1: "MaxId"}
)

// 设置 默认每页 10 条，页码 第 1 页
func SetDefaultPaginateOptions(in *PageOptions) *PageOptions {

	// set default pageSize ：
	if in.PageSize < 1 {
		in.PageSize = DefaultPageSize
	}

	// set default first page : 1
	if in.PageNumber < 1 {
		in.PageNumber = 1
	}

	return in
}

// 获取 select 参数
func GetSelectOptions(
	in *PageOptions, whereStr string, whereParams []interface{},
) (
	where string, params []interface{}, orderBy, offset interface{},
) {
	SetDefaultPaginateOptions(in)
	switch PaginateMode[in.PaginateMode] {
	case "MaxId":
		whereStr, whereParams, orderBy, offset = GetMaxIdSelectOptions(in, whereStr, whereParams)
	case "PageNumber":
		orderBy, offset = GetPageNumberSelectOptions(in)
	}

	return whereStr, whereParams, orderBy, offset
}

// PaginateMode ： PageNumber
func GetPageNumberSelectOptions(in *PageOptions) (orderBy, offset interface{}) {

	offset = in.PageSize * (in.PageNumber - 1)
	orderBy = in.OrderBy

	return
}

// PaginateMode ： MaxId ; default orderBy : id desc
func GetMaxIdSelectOptions(
	in *PageOptions, whereStr string, whereParams []interface{},
) (
	where string, params []interface{}, orderBy, offset interface{},
) {
	if IsOrderByAsc(in.OrderBy) {
		whereStr += " and id > ? "
		whereParams = append(whereParams, in.MaxId)
		orderBy = "id asc"
		offset = 0
	} else {
		if in.MaxId > 0 {
			whereStr += " and id < ? "
			whereParams = append(whereParams, in.MaxId)
		}
		offset = 0
		orderBy = "id desc"
	}

	return whereStr, whereParams, orderBy, offset
}

// 升序排序
func IsOrderByAsc(orderBy string) bool {
	return strings.HasSuffix(strings.ToLower(orderBy), " asc")
}

// set paginate data
func SetPaginateData(in *PageOptions, records, lastId int32) (paginate PageResult) {

	paginate.TotalRecords = records
	paginate.PageSize = in.PageSize
	paginate.PaginateMode = PaginateMode[in.PaginateMode]

	if paginate.TotalRecords%in.PageSize == 0 {
		paginate.TotalPages = paginate.TotalRecords / in.PageSize
	} else {
		paginate.TotalPages = paginate.TotalRecords/in.PageSize + 1
	}

	switch PaginateMode[in.PaginateMode] {
	case "PageNumber":
		paginate.PageNumber = in.PageNumber
	case "MaxId":
		if IsOrderByAsc(in.OrderBy) {
			paginate.MaxId = lastId
		} else {
			paginate.MaxId = lastId
		}
	}

	return
}
