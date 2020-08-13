package common

func Checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	DBHandle = 1000
	DBCache  = 1024
)
