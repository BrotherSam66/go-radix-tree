// Package radix
// @Title 基数树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package radix

import (
	"fmt"
	"go-radix-tree/radix/radixmodels"
)

// ShowTree 逐层显示这个树，
// @b 展示树的开始位置（可以不从root开始展示）
// @Author  https://github.com/BrotherSam66/
func ShowTree(n *radixmodels.RadixNode) {
	if n == nil { // 先判断，后定义变量，免得浪费内存
		fmt.Println("这个树/分支是空的")
		return
	}
	//fmt.Printf("%#v\n", n)

	// ShowTemp 数据。数据可能是nil。最多10层，每层最多1000 数据
	var showTemp [10][1000]*radixmodels.RadixNode
	//totalLevel := 0             // 总层数
	//nowLevel := 0               // 当前层数
	//nnn := Name
	nowColumn := 0 // 当前列
	fmt.Printf("\n展示树：(父|Child数)路径|string载荷|int载荷...\\child-1路径\\child-2路径...")
	//ShowOneNode(global.Root)
	//return
	showTemp[0][0] = n // 来的最高位指针

	for i := 1; i < len(showTemp); i++ { // 循环每一层
		fmt.Println("")              // 先来一个换行
		if showTemp[i-1][0] == nil { // 上一层全nil，结束显示
			return
		}
		nowColumn = 0 // 当前列
		for j := 0; j < len(showTemp[i-1]); j++ {
			if showTemp[i-1][j] == nil { // 本行遍历结束
				break
			}

			ShowOneNode(showTemp[i-1][j]) // 显示遍历到的上一行的这个节点。显示没换行

			// 在本行填写上一行所有节点的所有儿子
			for k := 0; k < len(showTemp[i-1][j].Child); k++ { // 其他的所有右腿
				showTemp[i][nowColumn] = showTemp[i-1][j].Child[k]
				nowColumn++
			}
		}
	}

}

// ShowOneNode 展示单个节点
// @Author  https://github.com/BrotherSam66/
func ShowOneNode(n *radixmodels.RadixNode) {
	if n == nil {
		fmt.Printf("()nil")
		return
	}
	fmt.Printf(" , ") // 左右分割

	// show父节点
	if n.Parent == nil {
		fmt.Printf("(nil|%d)", n.ChildNum)
	} else {
		fmt.Printf("(%s|%d)", string(n.Parent.Path), n.ChildNum)
	}

	// show本节点信息 路径|string载荷|int载荷...
	fmt.Printf("%s|%s|%v", string(n.Path), n.Payload, n.PayloadIntSlice)

	// show孩子 \child-1路径\child-2路径...
	for i := 0; i < n.ChildNum; i++ {
		fmt.Printf("\\%s", string(n.Child[i].Path))
	}
}
