package helper

// 对比两个不同int32数组
// compareOld : old->new, old 对比 new，在 new 中不存在
// compareNew : new->old, new 对比 old, 在 old 中不存在
func CompareSliceInt32(old, new *[]int32) (compareOld, compareNew []int32) {

	if len(*new) == 0 {
		compareOld = *old
		return
	}

	if len(*old) == 0 {
		compareNew = *new
		return
	}

	oldMap := make(map[int32]int)
	newMap := make(map[int32]int)

	for _, r := range *old {
		oldMap[r] = 1
	}

	for _, r := range *new {
		newMap[r] = 1
		// new 的 元素 不存在 old 中
		if _, ok := oldMap[r]; !ok {
			compareNew = append(compareNew, r)
		}
	}

	for _, r := range *old {
		// old 的 元素 不存在 new 中
		if _, ok := newMap[r]; !ok {
			compareOld = append(compareOld, r)
		}
	}

	return
}
