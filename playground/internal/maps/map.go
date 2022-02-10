package main

type Interface interface {
	Get() string
}

type Impl struct {
}

func (slf *Impl) Get() string {
	return "map"
}

func main() {
	var m = map[any][]int{}
	m[new(Impl)] = []int{1}
	m[Impl{}] = []int{1} // panic
}
