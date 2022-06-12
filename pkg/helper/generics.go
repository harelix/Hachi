package helper

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func MapKeys[K comparable, V any](m map[K]V, f func(K) K) map[K]V {
	for k, v := range m {
		m[f(k)] = v
	}
	return m
}

func CollectionFunc[K any](m []K, f func(K) bool) K {
	for _, v := range m {
		if f(v) {
			return v
		}
	}
	return *new(K)
}
