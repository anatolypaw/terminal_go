package entity

type Size struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

type InputEvent struct {
	Type     int
	Position Position
}
