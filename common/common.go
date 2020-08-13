package common

func Checkerr(err error)  {
	if err!=nil{
		panic(err)
	}
}