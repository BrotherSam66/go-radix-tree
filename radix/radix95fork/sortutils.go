package radix95fork

import "fmt"

// SelectSort 选择排序
func SelectSort(s *[]int) {
	length := len(*s)
	if length <= 1 { // 切片长度小于等于1，直接认为已经排好
		return
	}

	var maxValue int // 最大值
	var maxIndex int // 最大值出现的下标
	for i := 0; i < length; i++ {
		maxValue = (*s)[i]
		maxIndex = i
		for j := i + 1; j < length; j++ { // 从半截往后找最大的那个
			if maxValue < (*s)[j] {
				maxValue = (*s)[j]
				maxIndex = j
			}
		}
		if maxIndex != i { // 最大的不是第一个
			(*s)[i], (*s)[maxIndex] = (*s)[maxIndex], (*s)[i] // 交换
		}
		//fmt.Println(s)
	}
	return
}

// InsertSort 插入排序 (用链表有序插入感觉更快啊，测了，更慢)
func InsertSort(s *[]int) {
	length := len(*s)
	if length <= 1 { // 切片长度小于等于1，直接认为已经排好
		return
	}
	var tempValue int             // 拿出来准备插入的值
	var insertIndex int           // 锚点，准备插入的位置在锚点后面
	for i := 1; i < length; i++ { // 从下标1开始，下标0直接认为已经排好
		tempValue = (*s)[i]           // 准备拿来插入的值
		insertIndex = -1              // 表示原来顺序就对，不需要插入
		for j := i - 1; j > -1; j-- { // 从前半截排好的队列，倒过来找比我小的锚点，放着锚点后面
			if tempValue < (*s)[j] { // 如果我小于，则找到锚点。如果一直不大于，锚点=-1，就不用插入。
				break // 找到锚点，跳出
			}
			insertIndex = j // 不小于，潜在的锚点变换
		}
		if insertIndex != -1 { // 需要插入了，insertIndex后面就是插入锚点
			for k := i; k > insertIndex; k-- { // 循环，准备插入值的位置->锚点，倒序
				(*s)[k] = (*s)[k-1] // 逐个向右搬移一格
			}
			(*s)[insertIndex] = tempValue
		}
	}
	return
}

// QuickSort 快速排序
// left 我这次要处理的数据的最左端
// right 我这次要处理的数据的最右端
// s 完整的切片
func QuickSort(left int, right int, s *[]int) {
	// ①拿到数据，尾部数值作为支点值
	// ②小于支点值的，依次从左边填写；大于等于支点值的依次从右边填写
	// ③把支点值填写在最后剩余的位置。（不一定是中央，但是会是分界线）
	// ④左右新填写的的分别判断，哪边没结束，就递归调用下一层。
	//fmt.Println(left, right, s)

	if right-left <= 1 {
		if right-left < 1 {
			fmt.Println("出错了，right-left<1是不可能的")
			return
		}
		// 到这里肯定是==1，如果右小于左，就直接交换就好
		if (*s)[right] < (*s)[left] {
			(*s)[left], (*s)[right] = (*s)[right], (*s)[left] // 交换
		}
		return
	}

	// 到这里，right-left 长度>=2
	l := left  // 左边拟填充的位置
	r := right // 右边拟填充的位置
	//tempSlice := s[left : right-left] // 挑出来缓存的，准备重新布局的缓冲【不行，切片这个是引用类型】
	tempSlice := make([]int, right-left) // 挑出来缓存的，准备重新布局的缓冲。尾巴做中轴值了，不用缓存
	for i := 0; i < right-left; i++ {
		tempSlice[i] = (*s)[left+i]
	}
	pivotValue := (*s)[right]         // 中州，支点。用来居中的数值。这里取队尾
	for i := 0; i < right-left; i++ { // 不循环最后一个right是因为right作为中轴数据了
		if tempSlice[i] < pivotValue { // 小于中轴数据，贴左边放
			(*s)[l] = tempSlice[i]
			l++
		} else { // 大于等于中轴数据，贴左边放
			(*s)[r] = tempSlice[i]
			r--
		}
	}
	pivotIndex := l // 用来居中的数值的位置；一定是l和r之间有唯一的空挡
	(*s)[pivotIndex] = pivotValue

	if pivotIndex-left >= 2 { // 左边新填充 >= 2个。递归
		QuickSort(left, pivotIndex-1, s)
	}

	if right-pivotIndex >= 2 { // 右边新填充 >= 2个。递归
		QuickSort(pivotIndex+1, right, s)
	}
}
