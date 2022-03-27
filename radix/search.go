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

// Search 查找节点
// @key 键值
// @retNode 找到的节点指针（可能是适合插入的位置）
// @headKey  共同的头
// @tailKey  key，裁掉headKey剩余的
// @tailPath  path，裁掉headKey剩余的
// @Author https://github.com/BrotherSam66/
func (r *RadixNode) Search(key []byte) (retNode *RadixNode, headKey, tailKey, tailPath []byte, err error) {
	/*
		tempNode.Path 不可能为空
		key不可能为空，和EncodedPath必然前面若干相同(root.Path可能是空)
		[1]key=path（len(tailKey)==0 && len(tailPath)==0 ），==》完美找到，替换、补充值
		[2]path包含key（len(tailKey) == 0 && len(tailPath) != 0），==》就算找到。（回头在path末尾分叉）
		[3]path互不包含key（len(tailKey) != 0 && len(tailPath) != 0），==》就算找到。（回头在两个不同点分叉）
		[4]key包含path，len(tailKey) != 0  && len(tailPath) == 0
		[4.1]key包含path，==》如果tailKey首字节child==nil，就算找到。（回头在key末尾分叉）
		[4.2]key包含path，==》如果tailKey首字节child存在==》用 tailKey 向这个child递归
	*/
	if len(r.Path) == 0 && len(r.Child) == 0 {
		err = errors.New("这个树是空的！")
		fmt.Println(err.Error())
		return
	}

	// 处理 root 可能是0长度Path
	if len(r.Path) == 0 {
		// 必然是root.Path为空，但是有孩子==》[1]找到Key配对孩子==》递归;[2]找不到Key配对孩子==》在root加孩子;
		// 寻找哪个孩子和key配对
		childPoint, _, _ := FindChildPointInSlice(r, key[0])
		if childPoint < 0 { // 没有孩子配对，==》在root加孩子，root.Path为空，
			retNode = r
			headKey = []byte{}
			tailKey = key
			tailPath = []byte{}
		} else { // 有孩子配对，deep递归
			retNode, headKey, tailKey, tailPath, err = SearchDeep(r.Child[childPoint], key) // 深度搜索，需要递归
		}
		return
	}

	if key[0] != r.Path[0] {
		// key 和 root.Path 首字节不同，==》在root加孩子，root.Path为空，
		retNode = r
		headKey = []byte{}
		tailKey = key
		tailPath = []byte{} // 这时候，r.Path是有值的，故意回空，让上级识别
		return
	}

	retNode, headKey, tailKey, tailPath, err = SearchDeep(r, key) // 深度搜索，需要递归
	return

}

// SearchDeep 深度/递归查找节点
// @tempNode 从这个节点起步，向下找
// @key 键值
// @retNode 阶段性返回的节点（可能递归）
// @headKey  共同的头
// @tailKey  key，裁掉headKey剩余的
// @tailPath  path，裁掉headKey剩余的
// @Author  https://github.com/BrotherSam66/
func SearchDeep(tempNode *RadixNode, key []byte) (retNode *RadixNode, headKey, tailKey, tailPath []byte, err error) {
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

	if len(tempNode.Path) == 0 {
		err = errors.New("tempNode.Path 不可能为空")
		fmt.Println(err.Error())
		return
	}
	if len(key) == 0 {
		err = errors.New("key 不可能为空")
		fmt.Println(err.Error())
		return
	}
	if key[0] != tempNode.Path[0] {
		err = errors.New("key 和 tempNode.Path 首字节不同，不正常")
		fmt.Println(err.Error())
		return
	}
	// [4.2]key 包含 path ，==》如果有child，key首字母child存在==》砍短key，向这个child递归
	if len(key) > len(tempNode.Path) { // key 长，key做母串
		headKey, tailKey, tailPath, err = CompareByteSlice(&key, &tempNode.Path)
		if err != nil {
			return
		}
		if len(tailPath) == 0 { // 表示key 包含 path。有可能path还有child。只有有child==》才递归
			childPoint, _, _ := FindChildPointInSlice(tempNode, tailKey[0])
			if childPoint > -1 { // tempNode 有对应 key 的孩子
				retNode, headKey, tailKey, tailPath, err = SearchDeep(tempNode.Child[childPoint], tailKey) // 递归
				return
			}
		}
	} else { // key 短or等长，path做母串
		headKey, tailPath, tailKey, err = CompareByteSlice(&tempNode.Path, &key)
		if err != nil {
			return
		}
	}

	// 到这里，[1][2][3][4.1]情形，算是找到了
	retNode = tempNode

	return
}

// CompareByteSlice byte切片比较
// @a  母串
// @b  子串
// @headKey  共同的头
// @tailKeyA  母串，裁掉headKey剩余的
// @tailKeyB  子串，裁掉headKey剩余的
// @Author  https://github.com/BrotherSam66/
func CompareByteSlice(a, b *[]byte) (headKey, tailKeyA, tailKeyB []byte, err error) {
	if len(*a) < len(*b) {
		err = errors.New("A串不可以小于B串")
		fmt.Println(err.Error())
		return
	}
	if (a == nil) != (b == nil) {
		err = errors.New("A串、B串都不可以nil")
		fmt.Println(err.Error())
		return
	}
	i := 0
	for i = 0; i < len(*b); i++ {
		if (*b)[i] != (*a)[i] {
			break
		}
	}
	headKey = (*a)[0:i]
	tailKeyA = (*a)[i:]
	tailKeyB = (*b)[i:]
	return
}

// FindChildPointInSlice 找key值对应的child的点，
// @tempNode 被搜索的节点；
// @keyByte 搜索的关键值；
// @childPoint 找到的点位，=-1表示没找到；
// @middle 未必找到，但是如果插入，应该放入的点，
// @Author  https://github.com/BrotherSam66/
func FindChildPointInSlice(tempNode *RadixNode, keyByte byte) (childPoint int, middle int, err error) {
	childPoint = -1
	if tempNode == nil { // 空节点，不可能找到
		err = errors.New("被搜索节点是nil")
		fmt.Println(err.Error())
		return
	}
	if len(tempNode.Child) == 0 {
		return
	}

	left := 0
	right := len(tempNode.Child) - 1
	for {
		middle = (right-left)/2 + left
		if tempNode.Child[middle].Path[0] > keyByte { // 新加的小，向左
			if middle == left { // 已经是最左了，就算没找到
				return
			}
			right = middle - 1 // 右边界左移，递归
		} else if tempNode.Child[middle].Path[0] < keyByte { // 新加的大，向右
			if middle == right { // 已经是最右了，就算没找到
				middle++ // 将来插在最右
				return
			}
			left = middle + 1 // 左边界右移，递归
		} else { // 找到了
			return middle, middle, nil
		}
	}
}

// FindIntPointInSlice 找int值对应在intSlice的点，
// @intSlice 被搜索的；
// @inInt 搜索的关键值；
// @intPoint 找到的点位，=-1表示没找到；
// @middle 未必找到，但是如果插入，应该放入的点，
// @Author  https://github.com/BrotherSam66/
func FindIntPointInSlice(intSlice []int, inInt int) (intPoint int, middle int) {
	intPoint = -1
	left := 0
	right := len(intSlice) - 1
	for {
		middle = (right-left)/2 + left
		if intSlice[middle] > inInt { // 新加的小，向左
			if middle == left { // 已经是最左了，就算没找到
				return
			}
			right = middle - 1 // 右边界左移，递归
		} else if intSlice[middle] < inInt { // 新加的大，向右
			if middle == right { // 已经是最右了，就算没找到
				middle++ // 将来插在最右
				return
			}
			left = middle + 1 // 左边界右移，递归
		} else { // 找到了
			return middle, middle
		}
	}

}
