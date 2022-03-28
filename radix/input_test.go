// Package radix
// @Title 基数树测试
// @Description
// @Author  https://github.com/BrotherSam66/
// @Update
package radix

import (
	"reflect"
	"testing"
)

// 单挑测试 go test
// 测试代码覆盖率 go test -cover
// 多条测试 go test -v
// 基准测试。 在本目录下运行 go test -bench=InsertIntInSlice -benchmem

// 单条测试。 在本目录下运行 go test
func TestInsertIntInSlice(t *testing.T) {
	node := NewRadixNode(nil, []byte{}, "payload", 1)
	node.PayloadIntSlice = []int{1, 2, 3, 6, 7, 8}
	want := []int{1, 2, 3, 5, 6, 7, 8}
	node.InsertIntInSlice(5, 3)
	got := node.PayloadIntSlice
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want=:%v  got=:%v", want, got)
	}
}

// 多条测试。 在本目录下运行 go test -v
func TestInsertIntInSlice_1(t *testing.T) {
	type test struct {
		node     *RadixNode
		intSlice []int
		inInt    int
		inPoint  int
		want     []int
	}
	tests := map[string]test{
		"test_1": {NewRadixNode(nil, []byte{}, "payload", 1), []int{1, 2, 3, 6, 7, 8}, 5, 3, []int{1, 2, 3, 5, 6, 7, 8}},
		"test_2": {NewRadixNode(nil, []byte{}, "payload", 1), []int{1, 2, 3, 6, 7, 8}, 5, 3, []int{1, 2, 3, 5, 6, 7, 8}},
		"test_3": {NewRadixNode(nil, []byte{}, "payload", 1), []int{1, 2, 3, 6, 7, 8}, 5, 3, []int{1, 2, 3, 5, 6, 7, 8}},
		"test_4": {NewRadixNode(nil, []byte{}, "payload", 1), []int{1, 2, 3, 6, 7, 8}, 5, 3, []int{1, 2, 3, 5, 6, 7, 8}},
	}

	for k, tc := range tests {
		t.Run(k, func(t *testing.T) {
			tc.node.PayloadIntSlice = tc.intSlice
			tc.node.InsertIntInSlice(tc.inInt, tc.inPoint)
			if !reflect.DeepEqual(tc.want, tc.node.PayloadIntSlice) {
				t.Errorf("testName=:%s  want=:%v  got=:%v", k, tc.want, tc.node.PayloadIntSlice)
			}
		})
	}
}

// 基准测试。 在本目录下运行 go test -bench=InsertIntInSlice -benchmem
func BenchmarkInsertIntInSlice(b *testing.B) {
	node := NewRadixNode(nil, []byte{}, "payload", 1)
	for i := 0; i < b.N; i++ {
		node.PayloadIntSlice = []int{1, 2, 3, 6, 7, 8}
		node.InsertIntInSlice(5, 3)
	}
}
