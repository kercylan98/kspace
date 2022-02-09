package main

import "fmt"

func main() {
	var slice = make([]int, 5)
	fmt.Println(len(slice), cap(slice))
	slice = append(slice, []int{0, 1, 2, 3, 4}...)
	fmt.Println(len(slice), cap(slice))
	slice = append(slice, 5)
	fmt.Println(len(slice), cap(slice))
	slice = append(slice, []int{0, 1, 2, 3, 4, 6, 7, 8, 9, 10}...)
	fmt.Println(len(slice), cap(slice))

	var array = [5]int{1, 2, 3}
	fmt.Println(len(array), cap(array))
	fmt.Println()

	var sliceData = make([]int, 3, 10)
	sliceData[1] = 1
	sliceData[2] = 2
	var arrayData = [3]int{0, 1, 2}
	fmt.Println(sliceData, cap(sliceData))
	fmt.Println(arrayData)
	updateSlice(&sliceData)
	updateArray(arrayData)
	fmt.Println(sliceData, cap(sliceData))
	fmt.Println(arrayData)
}

// 如果传入的切片非指针类型，那么修改已有的数据生效，新增额外的数据无效
// 如果要使用append新增切片内容，必须使用指针类型
func updateSlice(data *[]int) {
	(*data)[0] = -1
	*data = append(*data, 99)
}

func updateArray(data [3]int) {
	data[0] = -1
}
