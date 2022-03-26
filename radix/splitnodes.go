package radix

import (
	"errors"
	"fmt"
	"go-radix-tree/radix/radixglobal"
	"go-radix-tree/radix/radixmodels"
)

// SplitNewNode 新节点分列，key砍短，新建 newDownNode 在下，tempNode在上，孩子跟下半截
// @key 键值
// @tempNode 找到的节点指针（可能是适合插入的位置）
// @headKey  共同的头
// @tailKey  key，裁掉headKey剩余的
// @tailPath  path，裁掉headKey剩余的
// @Author  https://github.com/BrotherSam66/
func SplitNewNode(tempNode *radixmodels.RadixNode, headKey, tailKey, tailPath []byte, payload string, payloadInt int) (err error) {
	// 新下位节点+新下位节点父亲+新下位节点内容
	newDownNode := radixmodels.NewRadixNode(tempNode, tailKey, payload, payloadInt)
	// 新下位节点孩子，没有

	// 老上位节点父亲，不变
	// 老上位节点内容，不变
	// 老上位节点孩子
	childPoint, insertPoint, _ := FindChildPointInSlice(tempNode, tailKey[0]) // 在newUpNode找到旧节点应该在的孩子的位置
	if childPoint < 0 {
		InsertChildInSlice(tempNode, newDownNode, insertPoint)
	} else {
		err = errors.New("tempNode，原来就包含了 newDownNode 这个孩子，这不科学")
		fmt.Println(err.Error())
		return
	}

	if payloadInt == 1 { // 打断点测试用
		return
	}

	return
}

// SplitOldNode 旧节点分列，path砍短，新建 newUpNode 在上，tempNode在下，孩子跟下半截
// @key 键值
// @tempNode 找到的节点指针（可能是适合插入的位置）
// @headKey  共同的头
// @tailKey  key，裁掉headKey剩余的
// @tailPath  path，裁掉headKey剩余的
// @Author  https://github.com/BrotherSam66/
func SplitOldNode(tempNode *radixmodels.RadixNode, headKey, tailKey, tailPath []byte, payload string, payloadInt int) (err error) {
	// 新上位节点+新上位节点父亲+新上位节点内容
	newUpNode := radixmodels.NewRadixNode(tempNode.Parent, headKey, payload, payloadInt)
	if newUpNode.Parent == nil { // 说明是root
		radixglobal.Root = newUpNode
	}
	// 新上位节点孩子
	if len(tailPath) == 0 { // 打断点位置
		fmt.Println("len(tailPath) == 0")
	}
	childPoint, insertPoint, _ := FindChildPointInSlice(newUpNode, tailPath[0]) // 在newUpNode找到旧节点应该在的孩子的位置
	if childPoint < 0 {
		InsertChildInSlice(newUpNode, tempNode, insertPoint)
	} else {
		err = errors.New("新上位节点newUpNode，原来就包含了tempNode这个孩子，这不科学")
		fmt.Println(err.Error())
		return
	}

	// 老下位节点父亲
	tempNode.Parent = newUpNode // 老节点降级
	// 老下位节点内容，
	tempNode.Path = tailPath // 老节点路径
	// 老下位节点孩子（不变）

	if payloadInt == 1 { // 打断点测试用
		return
	}

	return
}

// Split3Node 分裂成3个节点（路径分叉）
// @tempNode 旧节点
// @subKey 拟增加的节点编码路径
// @payload 载荷
// @payloadInt 数字载荷
// @author https://github.com/BrotherSam66/
func Split3Node(tempNode *radixmodels.RadixNode, headKey, tailKey, tailPath []byte, payload string, payloadInt int) (err error) {
	// 分裂成3个节点，[1]SplitOldNode分裂 newUpNode，payload随意；[2]针对newUpNode，SplitNewNode分裂newDownNode；[3]newUpNode的payload清空

	// [1]SplitOldNode分裂 newUpNode，payload随意；
	_ = SplitOldNode(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
	// [2]针对newUpNode，SplitNewNode分裂newDownNode；(tempNode.Parent就是newUpNode)
	_ = SplitNewNode(tempNode.Parent, headKey, tailKey, tailPath, payload, payloadInt)
	// [3]newUpNode的payload清空
	tempNode.Parent.Payload = ""
	tempNode.Parent.PayloadIntSlice = []int{}

	return
}
