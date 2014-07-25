package vutil

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}
