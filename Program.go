package main

func main() {
	Init()
	st := "test"
	it := int64(15)
	ft := "tesr"
	zt := int64(14)
	Add(&st, &st, &it)
	Add(&ft, &ft, &zt)
	t := "tes"
	values := Get(&t)
	for index := 0; index < len(values); index++ {
		println(values[index].key)
	}
}
