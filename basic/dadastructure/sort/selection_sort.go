package sort

/*
归并排序（merge sort）是一种基于分治策略的排序算法，包含图 11-10 所示的“划分”和“合并”阶段。

划分阶段：通过递归不断地将数组从中点处分开，将长数组的排序问题转换为短数组的排序问题。
合并阶段：当子数组长度为 1 时终止划分，开始合并，持续地将左右两个较短的有序数组合并为一个较长的有序数组，直至结束
*/
/*
快速排序（quick sort）是一种基于分治策略的排序算法，运行高效，应用广泛。
1. 划分哨兵，把数组不断划分到左边都小于base，右边大于base
2. 因为重复子问题，所以递归，改变数组的左右边界
快速排序的核心操作是“哨兵划分”，其目标是：选择数组中的某个元素作为“基准数”，将所有小于基准数的元素移到其左侧，而大于基准数的元素移到其右侧。具体来说，哨兵划分的流程如图 11-8 所示。
*/
/* 快速排序 */
func quickSort(nums []int, left, right int) {
	// 子数组长度为 1 时终止递归
	if left >= right {
		return
	}
	// 哨兵划分
	pivot := partition(nums, left, right)
	// 递归左子数组、右子数组
	quickSort(nums, left, pivot-1)
	quickSort(nums, pivot+1, right)
}
func partition(nums []int, left, right int) int {
	// 以 nums[left] 为基准数
	i, j := left, right
	for i < j {
		for i < j && nums[j] >= nums[left] {
			j-- // 从右向左找首个小于基准数的元素
		}
		for i < j && nums[j] <= nums[left] {
			i++
		}
		nums[i], nums[j] = nums[j], nums[i]
	}
	// 将基准数交换至两子数组的分界线
	nums[i], nums[left] = nums[left], nums[i]
	return i // 返回基准数的索引
}

/*
插入排序（insertion sort）是一种简单的排序算法，它的工作原理与手动整理一副牌的过程非常相似。
具体来说，我们在未排序区间选择一个基准元素，将该元素与其左侧已排序区间的元素逐一比较大小，并将该元素插入到正确的位置。
*/
func insertionSort(nums []int) {
	n := len(nums)
	for i := 1; i < n; i++ {
		base := nums[i]
		j := i - 1
		// 内循环：将 base 插入到已排序区间 [0, i-1] 中的正确位置
		for j >= 0 && nums[j] > base {
			nums[j+1] = nums[j]
			j--
		}
		nums[j+1] = base
	}
}

/*
冒泡排序（bubble sort）通过连续地比较与交换相邻元素实现排序。这个过程就像气泡从底部升到顶部一样，因此得名冒泡排序。
*/
func bubbleSort(nums []int) {
	// 外循环：未排序区间为 [0, i]
	for i := len(nums) - 1; i > 0; i-- {
		// 内循环：将未排序区间 [0, i] 中的最大元素交换至该区间的最右端
		for j := 0; j < i; j++ {
			if nums[j] > nums[j+1] {
				// 交换 nums[j] 与 nums[j + 1]
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
}

/*
选择排序（selection sort）的工作原理非常简单：开启一个循环，每轮从未排序区间选择最小的元素，将其放到已排序区间的末尾。
*/
func selectionSort(nums []int) {
	n := len(nums)
	// 外循环：未排序区间为 [i, n-1]
	for i := 0; i < n-1; i++ {
		// 内循环：找到未排序区间内的最小元素
		k := i
		for j := i + 1; j < n; j++ {
			if nums[j] < nums[k] {
				// 记录最小元素的索引
				k = j
			}
		}
		// 将该最小元素与未排序区间的首个元素交换
		nums[i], nums[k] = nums[k], nums[i]

	}
}
