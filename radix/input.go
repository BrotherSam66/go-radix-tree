// Package radix
// @Title 基数树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package radix

import (
	"errors"
	"fmt"
	"go-radix-tree/radix/radixglobal"
	"go-radix-tree/radix/radixmodels"
)

// Inputs 连续插入节点
// @author https://github.com/BrotherSam66/
func Inputs() {

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

		if payloadInt == 1 {
			payloadInt = 1
		}
		_ = Insert([]byte(key), "", payloadInt)
		ShowTree(radixglobal.Root)
		//if key == "-1" {
		//	return
		//}
		//if key == 0 {
		//	key = rand.Intn(radixglobal.MaxKey)
		//	fmt.Println(key)
		//}
		//if key > 1000 {
		//	if key > 1046 {
		//		fmt.Println("最大1046，否则溢出....")
		//		continue
		//	}
		//	endKey := key - 1000
		//	for i := 1; i <= endKey; i++ {
		//		//Insert(i, "")
		//	}
		//
		//	//ShowTree(radixglobal.Root)
		//	continue
		//}
		//if key > 99 || key < 1 {
		//	fmt.Println("必须是0~~99")
		//	continue
		//}
		//Insert(key, "")
		//ShowTree(radixglobal.Root)
	}
}

// Insert 加入节点
// @key 插入的键值
// @payload 插入的载荷值
// @author https://github.com/BrotherSam66/
func Insert(key []byte, payload string, payloadInt int) (err error) {
	if payload == "" {
		payload = string(key)
	}
	if radixglobal.Root == nil { // 原树为空树，新加入的转为根
		radixglobal.Root = radixmodels.NewRadixNode(nil, key, payload, payloadInt)
		return
	}

	// 从root开始查找附加的位置；tempNode=找到的节点。
	tempNode, headKey, tailKey, tailPath, err := Search(key)

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

	// [1]key=path==》完美找到，替换、补充值
	if len(tailKey) == 0 && len(tailPath) == 0 {
		err = PayloadModify(tempNode, payload, payloadInt)
		return
	}

	// [2]path包含key，==》path砍短，新建newUpNode在上，tempNode在下
	if len(tailKey) == 0 && len(tailPath) != 0 {
		_ = SplitOldNode(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
		return
	}

	// [3]path互不包含key，==》（回头在两个不同点分叉），上半截造一个全空纯粹分支节点，两个尾巴做两个叶子节点
	if len(tailKey) != 0 && len(tailPath) != 0 { // 互相不包含
		_ = Split3Node(tempNode, headKey, tailKey, tailPath, payload, payloadInt)
		return
	}

	// [4.1]key包含path，==》如果key首字母child==nil，（回头在path末尾分叉）
	if len(tailKey) != 0 && len(tailPath) == 0 { // key 包含 path
		childPoint, _, _ := FindChildPointInSlice(tempNode, tailKey[0]) // 在newUpNode找到旧节点应该在的孩子的位置
		if childPoint < 0 {                                             // key首字母child==nil
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
func PayloadModify(tempNode *radixmodels.RadixNode, payload string, payloadInt int) (err error) {
	intPoint, insertPoint := FindIntPointInSlice(tempNode.PayloadIntSlice, payloadInt)
	if intPoint < 0 { // 原节点不存在这个数值载荷，就插入
		tempNode.PayloadIntSlice = InsertIntInSlice(tempNode.PayloadIntSlice, payloadInt, insertPoint) // 插入
	}
	tempNode.Payload = payload
	//fmt.Println("tempNode.PayloadIntSlice", tempNode.PayloadIntSlice)
	return
}

// InsertChildInSlice intSlice指定位置插入inInt
// @intSlice 被插入的；
// @inInt 拟插入的值；
// @intPoint 拟插入的位置
// @return 返回的切片
// @Author  https://github.com/BrotherSam66/
func InsertChildInSlice(tempNode *radixmodels.RadixNode, inChild *radixmodels.RadixNode, insertPoint int) {
	childes := tempNode.Child
	childes = append(childes, inChild)                   // 切片扩展1个空间
	copy(childes[insertPoint+1:], childes[insertPoint:]) // a[i:]向后移动1个位置
	childes[insertPoint] = inChild                       // 设置新添加的元素
	tempNode.Child = childes
	tempNode.ChildNum++
	return
}

// InsertIntInSlice intSlice指定位置插入inInt
// @intSlice 被插入的；
// @inInt 拟插入的值；
// @intPoint 拟插入的位置
// @return 返回的切片
// @Author  https://github.com/BrotherSam66/
func InsertIntInSlice(intSlice []int, inInt int, insertPoint int) []int {
	intSlice = append(intSlice, 0)                         // 切片扩展1个空间
	copy(intSlice[insertPoint+1:], intSlice[insertPoint:]) // a[i:]向后移动1个位置
	intSlice[insertPoint] = inInt                          // 设置新添加的元素
	return intSlice
}
