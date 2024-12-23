package utils

func MapSlice[TSource any, TDest any](sourceData []TSource, selector func(TSource) TDest) []TDest {
    mapped := make([]TDest, len(sourceData))

    for i, item := range sourceData {
        mapped[i] = selector(item)
    }

    return mapped
}

func FilterSlice[TSource any](sourceData []TSource, predicate func(TSource) bool) []TSource {
    filtered := make([]TSource, 0, len(sourceData))

    for _, item := range sourceData {
        if predicate(item) {
            filtered = append(filtered, item)
        }
    }

    return filtered
}
