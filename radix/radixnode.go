// Package radix
// @Title 基数树模型
// @Author  https://github.com/BrotherSam66/
// @Update
package radix

// RadixNode 基数树节点的结构体
// @Author  https://github.com/BrotherSam66/
type RadixNode struct {
	Path            []byte       // 路径，字节的切片（否则无法处理汉字、二进制等的路径）
	Parent          *RadixNode   // 父节点
	MyCrc16         uint16       // 本节点的校验
	Child           []*RadixNode // 孩子节点，切片。按byte值定位，最大 256 叉
	ChildCrc16      []uint16     // 孩子节点校验，切片
	Payload         string       // 载荷
	PayloadIntSlice []int        // 载荷
	//ChildNum        int          // 孩子的个数

}

// NewRadixNode 构造函数
// @parent 父节点
// @path 路径
// @payload 载荷
// @payloadInt int载荷
// @Author https://github.com/BrotherSam66/
func NewRadixNode(parent *RadixNode, path []byte, payload string, payloadInt int) *RadixNode {
	return &RadixNode{
		Parent:          parent,
		Path:            path,
		Payload:         payload,
		PayloadIntSlice: []int{payloadInt},
	}
}
