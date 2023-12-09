package common

type Pair[R, S any] struct {
	First  R
	Second S
}

func MapFirst[R, S, T any](p Pair[R, S], f func(R) T) Pair[T, S] {
	return Pair[T, S]{
		First:  f(p.First),
		Second: p.Second,
	}
}

func MapSecond[R, S, T any](p Pair[R, S], f func(S) T) Pair[R, T] {
	return Pair[R, T]{
		First:  p.First,
		Second: f(p.Second),
	}
}

func MapSecondSlice[R, S, T any](p Pair[R, []S], f func(S) T) Pair[R, []T] {
	return Pair[R, []T]{
		First:  p.First,
		Second: MapSlice(p.Second, f),
	}
}

func MapSecondSliceSlice[R, S, T any](p Pair[R, [][]S], f func(S) T) Pair[R, [][]T] {
	return Pair[R, [][]T]{
		First: p.First,
		Second: MapSlice(p.Second, func(s []S) []T {
			return MapSlice(s, f)
		}),
	}
}

type Triple[R, S, T any] struct {
	First  R
	Second S
	Third  T
}
