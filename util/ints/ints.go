package ints

func Fill(data []int, value int) {
	for i := 0; i < len(data); i++ {
		data[i] = value
	}
}

func Index(data []int, value int) int {
	for i, v := range data {
		if v == value {
			return i
		}
	}
	return -1
}
