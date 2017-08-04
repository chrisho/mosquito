package paginate

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

	return in
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
 * if _, ok := requestParams["PagingByPrimary"]; ok {
 *	switch requestParams["SortById"].(int32) {
 *	case 1:
 *		sql = sql.Where(Table + ".id > ?", requestParams["PrimaryParam"])
 *		orderBy = SellPointLimitTable + ".id asc"
 *	default:
 *		if requestParams["PrimaryParam"].(int32) > 0 {
 *			sql = sql.Where(Table + ".id < ?", requestParams["PrimaryParam"])
 *		}
 *		orderBy = SellPointLimitTable + ".id desc"
 *	}
 * }
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
