// Package radix95fork
// @Title 基数树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package radix95fork

import (
	"errors"
	"fmt"
	"go-radix-tree/radix/radixglobal"
	"go-radix-tree/radix/radixmodels"
	"strings"
)

// Search 查找节点
// @key 键值
// @tempNode 找到的节点指针（可能是适合插入的位置）
// @isTarget 找到的是命中的节点
// @Author  https://github.com/BrotherSam66/
func Search(key string) (retNode *radixmodels.RadixNode, retKey string, err error) {
	if radixglobal.Root == nil {
		fmt.Println("这个树/分支是空的")
		return nil, "", errors.New("这个树是空的！")
	}

	retNode, retKey, err = SearchDeep(radixglobal.Root, key) // 深度搜索，需要递归
	return

}

func SearchDeep(tempNode *radixmodels.RadixNode, key string) (retNode *radixmodels.RadixNode, retKey string, err error) {

	/*
		tempNode.Path 不可能为空
		key不可能为空，和EncodedPath必然前面若干相同
		[1]key=path，==》完美找到
		[2]path包含key，==》就算找到。（回头在key末尾分叉）
		[3]path互不包含key，==》就算找到。（回头在两个不同点分叉）
		[4]key包含path，
		[4.1]key包含path，==》如果key首字母child==nil，就算找到（回头在path末尾分叉）
		[4.2]key包含path，==》如果有child，key首字母child存在==》砍短key，向这个child递归
	*/

	if tempNode.Path == "" {
		err = errors.New("tempNode.Path 不可能为空")
		fmt.Println(err.Error())
		return
	}
	if key == "" {
		err = errors.New("kry 不可能为空")
		fmt.Println(err.Error())
		return
	}
	if key[0:1] != tempNode.Path[0:1] {
		err = errors.New("kry 和 tempNode.Path 首字母不同，不正常")
		fmt.Println(err.Error())
		return
	}
	// [4.2]key 包含 path ，==》如果有child，key首字母child存在==》砍短key，向这个child递归
	subKey := "" // key 扣除path 后的串
	var r []rune
	if strings.Index(key, tempNode.Path) == 0 && len(key) > len(tempNode.Path) {
		subKey = key[len(tempNode.Path):] // key 扣除path 后的串
		r = []rune(subKey)                // 为了把subKey首字母转成int
	}

	if subKey != "" && tempNode.Child[int(r[0])-32] != nil {
		// [4.2]key 包含 path ，==》如果有child，key首字母child存在==》砍短key，向这个child递归.（-32是扣除ASCII表前面32个不可显示字符）
		retNode, retKey, err = SearchDeep(tempNode.Child[int(r[0])-32], subKey) // 递归
		return
	} else {
		// [1][2][3][4.1]情形，算是找到了
		retNode = tempNode
		retKey = key
		return
	}

}

// ChildPoint 定位孩子指针应该的位置。
// @str 编码路径字符串
// @Author  https://github.com/BrotherSam66/
func ChildPoint(str string) (childPoint int, err error) {
	childPoint = int(([]rune(str))[0]) - 32
	if childPoint < 0 || childPoint > 94 {
		err = errors.New(fmt.Sprintf("childPoint = %d,超出范围了", childPoint+32))
		fmt.Println(err.Error())
	}
	return
}

//// PredecessorOrSuccessor 找前驱or后继Key。比我稍小的最大Key，比我大的最小key
//// @key 键值
//// @avatar 找到的替身节点指针
//// @Author  https://github.com/BrotherSam66/
//func PredecessorOrSuccessor(tempNode *radixmodels.RadixNode, key int, isPredecessor bool) (avatar *radixmodels.RadixNode, err error) {
//	if tempNode.Child[0] == nil {
//		err = errors.New("已经是叶子了，不可以找前驱or后继")
//		fmt.Println("已经是叶子了，不可以找前驱or后继")
//		return
//	}
//
//	// 精准找到拟删除KEY的位置，deletePoint
//	deletePoint := 0
//	for deletePoint = 0; deletePoint < tempNode.KeyNum; deletePoint++ {
//		if tempNode.Key[deletePoint] == key { // 准确命中，只可能是新创建节点情形
//			break
//		}
//	}
//	if deletePoint >= tempNode.KeyNum {
//		fmt.Println("发生某种错误，找到KEY又不存在了！ ")
//		return
//	}
//
//	if isPredecessor { // 前驱
//		avatar = tempNode.Child[deletePoint] // 命中点左边的腿
//	} else { // 后继
//		avatar = tempNode.Child[deletePoint+1] // 命中点右边的腿
//	}
//
//	for { // 递归循环
//		if avatar.Child[0] == nil { // 到叶子，就算找到了，前驱Key=》尾巴，后继Key=头
//			return
//		}
//		// 下移一层
//		if isPredecessor { // 前驱
//			avatar = avatar.Child[avatar.KeyNum] // 不断向最右一条腿找
//		} else { // 后继
//			avatar = avatar.Child[0] // 不断向最左一条腿找
//		}
//	}
//}
