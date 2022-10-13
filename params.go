package splitwise

type Stringer interface {
	string
}

type Params[T Stringer] map[T]string
