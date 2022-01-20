package helper

type StringSlice []string

func (s StringSlice) Contains(str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func StringSliceContains(s []string, str string) bool {
	return StringSlice(s).Contains(str)
}

type IntSlice []int

func (s IntSlice) Contains(e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IntSliceContains(s []int, e int) bool {
	return IntSlice(s).Contains(e)
}
