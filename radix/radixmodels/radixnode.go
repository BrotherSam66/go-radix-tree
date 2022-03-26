// Package radixmodels
// @Title 基数树模型
// @Author  https://github.com/BrotherSam66/
// @Update
package radixmodels

import (
	"go-radix-tree/radix/radixconst"
)

// RadixNode 基数树节点的结构体
type RadixNode struct {
	Path            []byte       // 路径，字节的切片（否则无法处理汉字、二进制等的路径）
	Parent          *RadixNode   // 父节点
	MyCrc16         uint16       // 本节点的校验
	ChildNum        int          // 孩子的个数
	Child           []*RadixNode // 孩子节点，切片。按byte值定位，最大 256 叉
	ChildCrc16      []uint16     // 孩子节点校验，切片
	Payload         string       // 载荷
	PayloadIntSlice []int        // 载荷
}

// NewRadixNode 构造函数
func NewRadixNode(parent *RadixNode, path []byte, payload string, payloadInt int) *RadixNode {
	return &RadixNode{
		Parent:          parent,
		Path:            path,
		Payload:         payload,
		PayloadIntSlice: []int{payloadInt},
	}
}

// ===========下面的，已经废弃了===============

// Radix95Node 字典树节点的结构体
type Radix95Node struct {
	MyCrc16         uint16                     // 本节点的校验
	ChildNum        uint16                     // 孩子的个数
	Parent          *Radix95Node               // 父节点
	Path            string                     // 编码路径
	Child           [radixconst.M]*Radix95Node // 下级节点
	ChildCrc16      [radixconst.M]uint16       // 下级节点校验
	Payload         string                     // 载荷
	PayloadIntSlice []int                      // 载荷
}

// NewRadix95Node 构造函数
func NewRadix95Node(parent *Radix95Node, path string, payload string, payloadInt int) *Radix95Node {
	return &Radix95Node{
		Parent:          parent,
		Path:            path,
		Payload:         payload,
		PayloadIntSlice: []int{payloadInt},
	}
}

//func (n *RBTNode) ReplaceInfo(avatar *RBTNode) (err error) {
//	if avatar == nil {
//		return errors.New("老大，拟替换节点是nil，这活没法干啊！")
//	}
//	// n.IsRed  = avatar.IsRed // 不改颜色，不改连接指向
//	n.Key = avatar.Key
//	n.Label = avatar.Label
//	//n.Parent = avatar.Parent
//	//n.Left = avatar.Left
//	//n.Right = avatar.Right
//	return nil
//}
