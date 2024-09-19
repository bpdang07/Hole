package path

type Point struct {
	X float64
	Y float64
}

type Horizontal struct {
}

type path struct {
	head   string
	origin Point
}
