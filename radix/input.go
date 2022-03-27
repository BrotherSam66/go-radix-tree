// Package radix
// @Title 基数树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package radix

import (
	"errors"
	"fmt"
)

// Inputs 连续插入节点
// @author https://github.com/BrotherSam66/
func (r *RadixNode) Inputs() {

	for {
		var key string
		var payloadInt int
		fmt.Println("请输入KEY，按回车键(空按回车随机，-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == "-1" {
			return
		}

		fmt.Println("请输入payloadInt，按回车键：")
		_, _ = fmt.Scanln(&payloadInt)

		if payloadInt == 1 { // 打断点用
			payloadInt = 1
		}
		_ = r.Insert([]byte(key), "", payloadInt)
		r.ShowTree()
	}
}

// Insert 加入节点
// @key 插入的键值
// @payload 插入的载荷值
// @payloadInt 新(或添加的)数值载荷
// @author https://github.com/BrotherSam66/
func (r *RadixNode) Insert(key []byte, payload string, payloadInt int) (err error) {
	if len(key) == 0 {
		err = errors.New("没有输入key内容啊")
		fmt.Println(err.Error())
		return
	}
	if payload == "" {
		payload = string(key)
	}
	if len(r.Path) == 0 && len(r.Child) == 0 { // 原树为空树
		r.Path = key
		r.Payload = payload
		r.PayloadIntSlice = []int{payloadInt}
		return
	}

	// 从root开始查找附加的位置；tempNode=找到的节点。
	tempNode, headKey, tailKey, tailPath, err := r.Search(key)

	/*
		tempNode.Path 不可能为空
		key不可能为空，和EncodedPath必然前面若干相同
		[1]key=path（len(tailKey)==0 && len(tailPath)==0 ），==》完美找到，替换、补充值
		[2]path包含key（len(tailKey) == 0 && len(tailPath) != 0），==》就算找到。（回头在path末尾分叉）
		[3]path互不包含key（len(tailKey) != 0 && len(tailPath) != 0），==》就算找到。（回头在两个不同点分叉）
		[4]key包含path，len(tailKey) != 0  && len(tailPath) == 0
		[4.1]key包含path，==》如果tailKey首字节child==nil，就算找到。（回头在key末尾分叉）
		[4.2]key包含path，==》如果tailKey首字节child存在==》用 tailKey 向这个child递归
	*/

	// 在root加孩子，root.Path为空，
	if len(headKey) == 0 && len(tailKey) > 0 && len(tailPath) == 0 {
		if len(r.Path) == 0 { // r有其他孩子，在上面加一个key孩子就好
			childPoint, _, _ := FindChildPointInSlice(r, tailKey[0]) // 在newUpNode找到旧节点应该在的孩子的位置
			if childPoint < 0 {                                      // tailKey首字母child==nil（tempNode可能有其他孩子）
				_ = SplitNewNode(r, headKey, tailKey, r.Path, payload, payloadInt)
			} else { // tempNode.child 包含 tailKey这个分支，这是不可能的。
				err = errors.New("tempNode.child 包含 tailKey这个分支，这是不可能的。")
				fmt.Println(err.Error())
				return
			}
			return
		} else { //  r无其他孩子，3分叉，形成 空Path的root
			_ = r.Split3Node(r, headKey, tailKey, r.Path, payload, payloadInt)
			return
		}
	}

	// [1]key=path==》完美找到，替换、补充值
	if len(tailKey) == 0 && len(tailPath) == 0 {
		err = PayloadModify(tempNode, payload, payloadInt)
		return
	}

	// [2]path包含key，==》path砍短，新建newUpNode在上，tempNode在下
	if len(tailKey) == 0 && len(tailPath) != 0 {
		_, _ = r.SplitOldNode(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
		return
	}

	// [3]path互不包含key，==》（回头在两个不同点分叉），上半截造一个全空纯粹分支节点，两个尾巴做两个叶子节点
	if len(tailKey) != 0 && len(tailPath) != 0 { // 互相不包含
		_ = r.Split3Node(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
		return
	}

	// [4.1]key包含path，==》如果tailKey首字母child==nil，（回头在path末尾分叉）
	if len(tailKey) != 0 && len(tailPath) == 0 { // key 包含 path
		childPoint, _, _ := FindChildPointInSlice(tempNode, tailKey[0]) // 在newUpNode找到旧节点应该在的孩子的位置
		if childPoint < 0 {                                             // tailKey首字母child==nil（tempNode可能有其他孩子）
			_ = SplitNewNode(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
		} else { // tempNode.child 包含 tailKey这个分支，这是不可能的。
			err = errors.New("tempNode.child 包含 tailKey这个分支，这是不可能的。")
			fmt.Println(err.Error())
			return
		}
		return
	}
	return
}

// PayloadModify 修改节点的载荷数值
// @tempNode 被修改的节点
// @payload 新字符串载荷
// @payloadInt 新(或添加的)数值载荷
// @author https://github.com/BrotherSam66/
func PayloadModify(tempNode *RadixNode, payload string, payloadInt int) (err error) {
	if tempNode == nil {
		err = errors.New("PayloadModify出现 tempNode==nil")
		fmt.Println(err.Error())
		return
	}
	intPoint, insertPoint := FindIntPointInSlice(tempNode.PayloadIntSlice, payloadInt)
	if intPoint < 0 { // 原节点不存在这个数值载荷，就插入
		tempNode.InsertIntInSlice(payloadInt, insertPoint) // 插入
	}
	tempNode.Payload = payload
	//fmt.Println("tempNode.PayloadIntSlice", tempNode.PayloadIntSlice)
	return
}

// InsertChildInSlice 节点指定位置插入孩子
// @tempNode 被插入的节点
// @inChild 拟插入的孩子；
// @intPoint 拟插入的位置
// @Author  https://github.com/BrotherSam66/
func InsertChildInSlice(tempNode *RadixNode, inChild *RadixNode, insertPoint int) {
	// Slice 是引用类型，必须逐个元素搬移
	tempNode.Child = append(tempNode.Child, inChild) // 只是扩容
	// 逐个向后搬移。尾开始，结束insertPoint+1，刚好把搬走
	for i := len(tempNode.Child) - 1; i > insertPoint; i-- {
		tempNode.Child[i] = tempNode.Child[i-1]
	}
	tempNode.Child[insertPoint] = inChild
	//tempNode.ChildNum++
	return
}

// InsertIntInSlice 节点指定位置插入inInt
// @n 被插入的节点
// @inInt 拟插入的int；
// @intPoint 拟插入的位置
// @Author  https://github.com/BrotherSam66/
func (n *RadixNode) InsertIntInSlice(inInt int, insertPoint int) {
	// Slice 是引用类型，必须逐个元素搬移
	n.PayloadIntSlice = append(n.PayloadIntSlice, inInt) // 只是扩容
	// 逐个向后搬移。尾开始，结束insertPoint+1，刚好把搬走
	for i := len(n.PayloadIntSlice) - 1; i > insertPoint; i-- {
		n.PayloadIntSlice[i] = n.PayloadIntSlice[i-1]
	}
	n.PayloadIntSlice[insertPoint] = inInt
}
