package common

func Compose2[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

func Compose3[A, B, C, D any](f func(A) B, g func(B) C, h func(C) D) func(A) D {
	return func(a A) D {
		return h(g(f(a)))
	}
}

func Compose4[A, B, C, D, E any](f func(A) B, g func(B) C, h func(C) D, i func(D) E) func(A) E {
	return func(a A) E {
		return i(h(g(f(a))))
	}
}

func Compose5[A, B, C, D, E, F any](f func(A) B, g func(B) C, h func(C) D, i func(D) E, j func(E) F) func(A) F {
	return func(a A) F {
		return j(i(h(g(f(a)))))
	}
}

func Compose6[A, B, C, D, E, F, G any](f func(A) B, g func(B) C, h func(C) D, i func(D) E, j func(E) F, k func(F) G) func(A) G {
	return func(a A) G {
		return k(j(i(h(g(f(a))))))
	}
}
