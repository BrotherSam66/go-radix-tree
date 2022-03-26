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

//// InsertOneNode 插入一个节点，可能要递归
//// @n 被插入的节点
//// @insertNode 拟插入的节点，要么①新节点，儿子均nil；要么②下层满员把中间节点挤上来（上来前key放在Key[1]，把下层分裂，作为我的2个儿子，放在Child[0,1]）
//// @author https://github.com/BrotherSam66/
//func InsertOneNode(n *btreemodels.BTreeNode, insertNode *btreemodels.BTreeNode) (err error) {
//	// 寻找插入的位置，拟插入放在这个点前面
//	keyPoint := 0
//	for keyPoint = 0; keyPoint < n.KeyNum; keyPoint++ {
//		if insertNode.Key[0] == n.Key[keyPoint] { // 准确命中，只可能是新创建节点
//			n.Payload = insertNode.Payload
//			return
//		} else if insertNode.Key[0] < n.Key[keyPoint] { // 说明已经找过头了,结束本节点循环，key插在i前面
//			break
//		}
//		// 到这里：可能①会向后找；可能②KeyNum循环结束，得到的i是最右key的右边，拟插入key本组最大。
//	}
//
//	// 到这里：i表示了拟插入key的位置。insertNode可能是不带孩子的新创建节点，也可能是下层挤上来的带2个孩子的节点(不会凭空上来，有一条腿是要替换原来的父节点的，我们指定用左腿)
//	// 强行插入，无论是否满员，溢出的在Tail里
//	keyTail, payloadTail, childTail, _ := InsertOneKey(n, insertNode, keyPoint)
//	// 分析本节点是否需要裂变
//	if n.KeyNum < globalconst.M-1 { // 被插入节点不满员，不用递归
//		n.KeyNum++
//		return
//	}
//
//	// 到这里，被插入节点满员，就需要分裂，需要递归了
//	// 开始对本节点分裂，分裂成3个，升起的是M/2位置的
//	upNode, isUpRoot, _ := SplitTo3Node(n, keyTail, payloadTail, childTail)
//
//	// 这里只是把中间节点升起来，拟插入下一级，带着两条腿，进入下一层递归。（如果本节点是root，升起来的就是新root就结束）
//	if isUpRoot { // 说明升起来的是单root
//		n.Parent = upNode    // 左儿子重新认爹
//		radixglobal.Root = upNode // 重新指定根节点n.Parent = {*go-b-tree-bplus-tree/btreemodels.BTreeNode | 0xc000120500}
//		return
//	} else { // 不是root升起来的。递归...
//		tempNode := n.Parent              // 原来被插入的节点的爹作为新的被插入的节点，拿来递归的
//		upNode.Child[0].Parent = tempNode // 上升节点的两个儿子指向上升节点拟插入的节点
//		upNode.Child[1].Parent = tempNode // 上升节点的两个儿子指向上升节点拟插入的节点
//		//n.Parent = upNode                   // 原来被插入的节点当up节点的左儿子
//		_ = InsertOneNode(tempNode, upNode) // 递归
//		return
//	}
//	// 不可能到这里
//}
//
//// InsertOneKey 插入一个Key，满了也插，溢出在Tail里
//// @n 被插入节点
//// @insertNode 拟插入节点
//// @insertPoint 拟插入位置，新入的占用这个位置
//// @keyTail 准备承载Key数组最后一个元素
//// @ChildTail  准备承载Child数组最后一个元素
//// @payloadTail 准备承载payload数组最后一个元素
//// @author https://github.com/BrotherSam66/
///*
// *假设：5阶，最大4个KEY、最小2个KEY，孩子数=KEY数+1，(65)是从60|70中间原来指向节点分裂升上来的
// *  (20|30  |              80)   |  (20|30  |              80)    |  (20|30  |      60|       80)   |
// *  /   \    \                \  |  /   \    \                \   |  /   \    \        \         \  |
// *(?1)(?2) (40|50|60   |70)  (?3)|(?1)(?2) (40|50|60|65 | 70) (?3)|(?1)(?2) (40|50)     (65|70) (?3)|
// *         /   \  \        \     |         /   \  \  \   \   \    |         /   \  \    /   \  \    |
// *       (?4)(?5)(?6) (65) (?7)  |       (?4)(?5)(?6)(?8)(?9)(?7) |       (?4)(?5)(?6) (?8)(?9)(?7) |
// *                    / \        |                                |                                 |
// *                  (?8)(?9)     |                                |                                 |
// *(?8)是(65)原归属节点左半部分，原来就和60|70指针勾连，
// *(?9)是(65)原归属节点右半部分，是新分裂出来的。
// */
//func InsertOneKey(n *btreemodels.BTreeNode, insertNode *btreemodels.BTreeNode, insertPoint int) (keyTail int, payloadTail string, childTail *btreemodels.BTreeNode, err error) {
//	keyTail = n.Key[globalconst.M-2]         // 数组最后一个元素
//	payloadTail = n.Payload[globalconst.M-2] // 数组最后一个元素
//	childTail = n.Child[globalconst.M-1]     // 数组最后一个元素
//
//	// 把往后挤走的KEY处理完
//	// 例如 globalconst.M-1=9最大9个key；n.KeyNum=6目前6个key；①keyPoint=3表示拟插入要在3这个位置，②keyPoint=6表示拟插入最大
//	//for j := n.KeyNum; j > insertPoint; j-- { // 例如①KeyNum=4，insertPoint=0，j=3~1；②KeyNum=4，insertPoint=1，j=3~2
//	for j := globalconst.M - 2; j > insertPoint; j-- { // 咬死从数组最后一个元素倒其，可能浪费些算力
//		n.Key[j] = n.Key[j-1]
//		n.Payload[j] = n.Payload[j-1]
//		n.Child[j+1] = n.Child[j] // 搬移的是每个Key的右腿
//	}
//
//	// 把拟插入节点放进来
//	// 升上来的节点(不会凭空上来，有一条腿是要替换原来的父节点的，我们指定共享左腿，但是插入0位指定共享右腿)。下面有一句是废话
//	if insertPoint > globalconst.M-2 { // 插入的是在溢出的尾巴，实际上插入的才是溢出
//		keyTail = insertNode.Key[0]         // 溢出的key
//		payloadTail = insertNode.Payload[0] // 数组最后一个元素
//		childTail = insertNode.Child[1]     // 溢出的右腿。（左腿insertNode上来前已经确保和n的最右腿取值一样了）
//	} else { // 插入的不是在尾巴，真实插入
//		n.Key[insertPoint] = insertNode.Key[0]
//		n.Payload[insertPoint] = insertNode.Payload[0]
//		n.Child[insertPoint+1] = insertNode.Child[1] // 右腿
//		// todo 还要考虑插入的左右腿
//		if insertPoint == 0 {
//			n.Child[0] = insertNode.Child[0] // 左腿
//		}
//	}
//	return
//}
