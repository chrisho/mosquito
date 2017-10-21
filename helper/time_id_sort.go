package helper

// 默认长度为19为，time(10位）* 1E9(9位） + id， 修改长度，请修改改1E9...
func TimeIdSort(time, id int32) int64 {
	if time <= 0 {
		return 0
	}
	return int64(time) * 1E9 + int64(id)
}
