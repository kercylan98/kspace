package main

import "fmt"

type EmptyStruct struct{}

type CommonStruct struct {
	ID string
}

func main() {
	var emptyStruct = EmptyStruct{}
	fmt.Println(EmptyStruct{} == EmptyStruct{}, true)
	fmt.Println(&EmptyStruct{} == &EmptyStruct{}, false)
	fmt.Println(*(&EmptyStruct{}) == EmptyStruct{}, true)
	fmt.Println(&EmptyStruct{} == &(EmptyStruct{}), false)
	fmt.Println(emptyStruct == emptyStruct, true)
	fmt.Println(&emptyStruct == &emptyStruct, true)

	var commonStruct = CommonStruct{ID: "0"}
	var commonStructA = CommonStruct{ID: "A"}
	var commonStructB = CommonStruct{ID: "B"}
	fmt.Println(CommonStruct{ID: "1"} == CommonStruct{ID: "1"}, true)
	fmt.Println(CommonStruct{ID: "1"} == CommonStruct{ID: "2"}, false)
	fmt.Println(&CommonStruct{ID: "1"} == &CommonStruct{ID: "1"}, false)
	fmt.Println((&CommonStruct{ID: "1"}).ID == (&CommonStruct{ID: "1"}).ID, true)
	fmt.Println(commonStruct == commonStruct, true)
	fmt.Println(&commonStruct == &commonStruct, true)
	fmt.Println(&commonStructA == &commonStructB, false)
	fmt.Println(commonStructA == commonStructB, false)
}
