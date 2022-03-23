package radixmodels

import (
	"go-radix-tree/radix/radixconst"
)

// RadixNode 字典树节点的结构体
type RadixNode struct {
	MyCrc16     uint16                   // 本节点的校验
	EncodedPath string                   // 编码路径
	Child       [radixconst.M]*RadixNode // 下级节点
	ChildCrc16  [radixconst.M]uint16     // 下级节点校验
	Payload     string                   // 载荷
}

// NewRadixNode 构造函数
func NewRadixNode(myCrc16 uint16, encodedPath string, payload string) *RadixNode {
	return &RadixNode{
		MyCrc16:     myCrc16,
		EncodedPath: encodedPath,
		Payload:     payload,
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
