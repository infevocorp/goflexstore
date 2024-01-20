package gormstore

func defaultValue[T comparable](val T, defaultVal T) T {
	if val == (*new(T)) {
		return defaultVal
	}

	return val
}
