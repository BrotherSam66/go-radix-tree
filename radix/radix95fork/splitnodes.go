package radix95fork

import (
	"errors"
	"fmt"
	"go-radix-tree/radix/radixglobal"
	"go-radix-tree/radix/radixmodels"
)

func SplitOldNode(tempNode *radixmodels.RadixNode, subKey string, payload string, payloadInt int) (err error) {
	// 分裂旧节点，上半截放新内容，下半截放就内容，孩子跟下半截

	// 编码路径
	//encodedPathHead := tempNode.Path[0:len(subKey)]
	stringTailOld := tempNode.Path[len(subKey):] // 路由尾部，准备扣掉首字母给下位节点

	// 新上位节点+新上位节点父亲+新上位节点内容
	newNodeUp := radixmodels.NewRadixNode(tempNode.Parent, subKey, payload, payloadInt) // 新节点，准备在上位
	if newNodeUp.Parent == nil {                                                        // 说明是root
		radixglobal.Root = newNodeUp
	}
	// 新上位节点孩子
	r := []rune(stringTailOld)
	newNodeUp.Child[int(r[0])-32] = tempNode // 找对child位置，指向老节点
	newNodeUp.ChildNum++

	// 老下位节点父亲
	tempNode.Parent = newNodeUp // 老节点降级
	// 老下位节点内容，
	tempNode.Path = stringTailOld[1:] // 老节点编码路径，扣掉首字母
	// 老下位节点孩子（不变）

	if payloadInt == 1 { // 打断点测试用
		return
	}

	return
}

// Split3Node 分裂成3个节点（编码路径分叉）
// @tempNode 旧节点
// @subKey 拟增加的节点编码路径
// @payload 载荷
// @payloadInt 数字载荷
// @author https://github.com/BrotherSam66/
func Split3Node(tempNode *radixmodels.RadixNode, subKey string, payload string, payloadInt int) (err error) {
	// 分裂成3个节点，上半截造一个全空纯粹分支节点，两个尾巴做两个叶子节点

	// 编码路径
	stringHead, stringTailOld, stringTailNew, _ := Split3String(tempNode.Path, subKey)

	// 新上位节点+新上位节点父亲+新上位节点内容
	newNodeUp := radixmodels.NewRadixNode(tempNode.Parent, stringHead, "", -1) // 新节点，准备在上位
	if newNodeUp.Parent == nil {                                               // 说明是root
		radixglobal.Root = newNodeUp
	}
	newNodeUp.PayloadIntSlice = []int{}
	// 新上位节点孩子-old
	r := []rune(stringTailOld)
	newNodeUp.Child[int(r[0])-32] = tempNode // 找对child位置，指向老节点
	newNodeUp.ChildNum++
	// 新上位节点孩子-new（要等一等）

	// 新下位节点+新下位节点父亲+新下位节点内容(编码路径stringTailNew扣掉首字母)
	newNodeDown := radixmodels.NewRadixNode(newNodeUp, stringTailNew[1:], payload, payloadInt) // 新节点，准备在上位

	// 新上位节点孩子-new(现在搞)
	r = []rune(stringTailNew)
	newNodeUp.Child[int(r[0])-32] = newNodeDown // 找对child位置，指向老节点
	newNodeUp.ChildNum++
	// 新下位节点孩子（没有）

	// 老下位节点父亲
	tempNode.Parent = newNodeUp // 老节点降级
	// 老下位节点内容，
	tempNode.Path = stringTailOld[1:] // 老节点编码路径，扣掉首字母
	// 老下位节点孩子（不变）

	if payloadInt == 1 { // 断电测试用
		return
	}

	return
}

// Split3String 比对2个string，分解出相同和不同的部分
// @encodedPath 原始编码路径 待比对string 1
// @subKey 拟增加的节点编码路径 待比对string 2
// @stringHead 共同的头
// @stringTailOld 原始编码路径剩下的尾巴
// @stringTailNew 拟增加的节点编码路径剩下的尾巴
// @author https://github.com/BrotherSam66/
func Split3String(encodedPath, subKey string) (stringHead, stringTailOld, stringTailNew string, err error) {
	rOld := []rune(encodedPath)
	rNew := []rune(subKey)
	i := 0
	for i = 0; i < len(rOld); i++ {
		if i >= len(rNew)-1 {
			err = errors.New("错误，encodedPath 包含 subKey")
			fmt.Println(err.Error())
			return
		}
		if rOld[i] != rNew[i] { // 找到分叉点了
			if i >= len(rOld)-1 {
				err = errors.New("错误，subKey 包含 encodedPath")
				fmt.Println(err.Error())
				return
			}
			if i == 0 {
				err = errors.New("错误，subKey 和 encodedPath没有共同的头部")
				fmt.Println(err.Error())
				return
			}
			break
		}
	}
	stringHead = subKey[0:i]
	stringTailOld = encodedPath[i:]
	stringTailNew = subKey[i:]
	return
}
