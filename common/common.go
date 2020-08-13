package common

func Checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	DBHandle = 500
	DBCache  = 512
)
