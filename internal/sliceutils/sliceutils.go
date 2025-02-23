package sliceutils

func Deduplicate[C comparable](items []C) []C {
	seen := make(map[C]struct{}, max(16, len(items)/16))
	unique := make([]C, 0, len(items))

	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		unique = append(unique, item)
	}
	return unique
}
