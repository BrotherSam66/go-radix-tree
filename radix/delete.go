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

/*
[1]目标是叶子 or 分支(肯定有payload)带单孩，删除payload。==》转向[1.2]
[1.1]删除后，父亲只剩下多个孩子 or 父亲有payload。==》结束
[1.2]删除后，父亲只剩下1个孩子 and 父亲无payload，==》把兄弟向上合并给父亲==》用父亲递归
[2]目标是分支(肯定有payload)带多孩。删除payload。==》结束
*/

// Deletes 连续删除节点
// @author https://github.com/BrotherSam66/
func Deletes() {

	for {
		var key string
		fmt.Println("请输入KEY，按回车键(空按回车随机，-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == "-1" {
			return
		}

		_ = Delete([]byte(key))
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

// Delete 删除节点
// @key 键值
// @author https://github.com/BrotherSam66/
func Delete(key []byte) (err error) {

	if radixglobal.Root == nil { // 原树为空树，新加入的转为根
		err = errors.New("这是一颗空树")
		fmt.Println(err.Error())
		return
	}

	// 从root开始查找附加的位置；tempNode=找到的节点。必须完美找到
	tempNode, _, tailKey, tailPath, err := Search(key)

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
	if len(tailKey) != 0 || len(tailPath) != 0 {
		err = errors.New("没找到这个节点")
		fmt.Println(err.Error())
		return
	}
	if tempNode.Payload == "" && len(tempNode.PayloadIntSlice) == 0 && len(tempNode.Child) > 1 {
		err = errors.New("找到的是分支节点没有payload，无从删除")
		fmt.Println(err.Error())
		return
	}

	// 找到是可能带payload+有孩孩分支（单孩or多孩）==》递归，让递归去处理
	if len(tempNode.Child) > 0 {
		tempNode.Payload = ""
		tempNode.PayloadIntSlice = []int{}

		//递归
		_ = TryMergeParentAndSan(tempNode, tempNode.Child[0])
		return
	}

	// 找到是叶子==》删除叶子==》可能递归
	if len(tempNode.Child) == 0 {

		if tempNode.Parent == nil { // 我是叶子，也是是根节点.==》删空==》结束
			radixglobal.Root = nil
			return
		}
		parent := tempNode.Parent
		// 找到tempNode在san.child的位置childPoint
		childPoint, _, _ := FindChildPointInSlice(parent, tempNode.Path[0])
		if childPoint < 0 {
			err = errors.New("上级节点parent.Child里找不到tempNode的位置，无从删除")
			fmt.Println(err.Error())
			return
		}
		// 删除parent.child[childPoint].相当于删除了这个叶子
		copy(parent.Child[childPoint+1:], parent.Child[:childPoint-1]) // a[i:]向后移动1个位置

		//递归.parent.Parent 可能= nil
		_ = TryMergeParentAndSan(parent.Parent, parent)
		return

	}
	return
}

// TryMergeParentAndSan 尝试合并父子节点
// @parent  父节点
// @san  子节点
// @Author  https://github.com/BrotherSam66/
func TryMergeParentAndSan(parent, san *radixmodels.RadixNode) (err error) {
	// parent ==root ？？？？
	if parent == nil { // san 是 root
		if san.Payload == "" && len(san.PayloadIntSlice) == 0 && len(san.Child) == 0 { // 无payload+无孩子
			radixglobal.Root = nil
		}
		return
	}

	// parent 多孩==》结束
	if len(parent.Child) > 1 {
		return
	}

	// parent 有payload + san有payload ==》结束
	if (parent.Payload != "" || len(parent.PayloadIntSlice) > 0) &&
		(san.Payload != "" || len(san.PayloadIntSlice) > 0) {
		return
	}

	// (parent 单孩+无payload) ==》san合并到parent==》递归
	if parent.Payload == "" && len(parent.PayloadIntSlice) == 0 {
		_ = MergeParentAndSan(parent, san)
		_ = TryMergeParentAndSan(parent.Parent, parent) // 递归
		return
	}

	fmt.Println("TryMergeParentAndSan,不应该走到这里啊，也许是递归结束了")
	return
}

// MergeParentAndSan 真实合并父子节点
// @parent  父节点，合并后的节点，由parent代表
// @san  子节点
// @Author  https://github.com/BrotherSam66/
func MergeParentAndSan(parent, san *radixmodels.RadixNode) (err error) {
	// 这里 parent无payload + 有唯一的 san，
	copy(parent.Path, san.Path)                       // 路劲
	parent.Payload = san.Payload                      // payload
	copy(parent.PayloadIntSlice, san.PayloadIntSlice) // payload
	parent.ChildNum = san.ChildNum                    // childNum
	parent.Child = parent.Child[0:0]                  // 删除唯一的下联键
	copy(parent.Child, san.Child)                     // child
	// san 的孩子的上联
	for i := 0; i < len(parent.Child); i++ {
		parent.Child[i].Parent = parent
	}

	return

}
