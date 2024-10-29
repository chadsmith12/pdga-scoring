package utils

func MapSlice[TSource any, TDest any](sourceData []TSource, selector func(TSource) TDest) []TDest {
    mapped := make([]TDest, len(sourceData))

    for i, item := range sourceData {
        mapped[i] = selector(item)
    }

    return mapped
}
