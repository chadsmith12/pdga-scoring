package utils_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/utils"
)

func TestMapSlices(t *testing.T) {
    source := []int{1, 2, 3, 4, 5}    
    doubled := utils.MapSlice(source, func(i int) int {
       return i * 2 
    })
    
    expected := []int{2, 4, 6, 8, 10}
    assertSlicesAreSame(t, expected, doubled)
}

func TestFilterSlices(t *testing.T) {
    source := []int{1, 2, 3, 4, 5, 6}
    onlyEven := utils.FilterSlice(source, func(item int) bool {
        return item % 2 == 0
    })

    expected := []int{2, 4, 6}
    assertSlicesAreSame(t, expected, onlyEven)
}

func assertSlicesAreSame(t *testing.T, expected []int, actual []int) {
    if len(expected) != len(actual) {
        t.Fatalf("len(expected) != len(actual): expected: %d; got: %d", len(expected), len(actual))
    }

    for i, item := range expected {
        if item != expected[i] {
            t.Fatalf("expected[%d] != actual[%d]: expected: %d, got: %d", i, i, item, expected[i])
        }
    }
}
