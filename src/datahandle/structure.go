package datahandle

// list []string 去重
func RemoveRepStr(slc []string) []string {
	var removeRepByMap = func(slc []string) []string {
		result := make([]string, 0)
		tempMap := map[string]byte{}
		for _, e := range slc {
			l := len(tempMap)
			tempMap[e] = 0
			if len(tempMap) != l {
				result = append(result, e)
			}
		}
		return result
	}

	var removeRepByLoop = func(slc []string) []string {
		result := make([]string, 0)
		for i := range slc {
			flag := true
			for j := range result {
				if slc[i] == result[j] {
					flag = false
					break
				}
			}
			if flag {
				result = append(result, slc[i])
			}
		}
		return result
	}

	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return removeRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return removeRepByMap(slc)
	}
}
