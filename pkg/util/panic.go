package util

func PanicIfNotNull(e interface{}) {
	if e != nil {
		panic(e)
	}
}

func PanicIfNull(e interface{}) {
	if e == nil {
		panic(e)
	}
}
