package cmd

type SliceMap[I, O any] func(I) O

func SliceMapping[I, O any](inputs []I, from SliceMap[I, O]) (output []O) {
	output = make([]O, len(inputs))
	for i, input := range inputs {
		output[i] = from(input)
	}
	return
}

type MapToSlice[K, V, O any] func(K, V) O

func MapToSliceMapping[K comparable, V, O any](inputs map[K]V, from MapToSlice[K, V, O]) (output []O) {
	output = make([]O, len(inputs))
	position := 0
	for key, value := range inputs {
		output[position] = from(key, value)
		position++
	}
	return
}

func ReturnKey[K, V any](key K, ignored V) K {
	return key
}
