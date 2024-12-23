package results

import "github.com/chadsmith12/pdga-scoring/pkgs/pulse"

type ListResult[T any] struct {
    Data []T `json:"data"`
    Size int `json:"size"`
}

func List[T any](data []T) pulse.JsonResultWriter {
    if (len(data) == 0) {
        data = []T{}
    }

    result := ListResult[T] {
        Data: data,
        Size: len(data),
    }

    return pulse.JsonResult(result)
}

