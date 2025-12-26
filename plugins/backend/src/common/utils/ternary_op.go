package utils

func TernaryOp[T any](cond bool, trueCase, falseCase T) T {
	if cond {
		return trueCase
	} else {
		return falseCase
	}
}
