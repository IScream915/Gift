package lottery

import "fmt"

// Lottery 抽奖函数
// 入参: 抽奖物品的库存量
// 出参: 抽中的物品的编号
func Lottery(items []int) int {
	// items 是参与抽奖的物品个数切片, 例如:
	// [2, 2, 4, 7]		item: 参与抽奖的每个物品的个数
	// [0, 1, 2, 3] 	index: 抽奖物品的编号
	// [2, 4, 8, 15]	acc: 抽奖物品的累积个数, 0号物品有2-0=2个, 1号物品有4-2=2个, 2号物品有8-4=4个, 3号物品有15-8=7个
	sum := 0              // sum用于保存参与抽奖物品的总个数
	acc := make([]int, 0) // acc用于保存抽奖物品的累积个数

	for _, item := range items {
		sum += item
		acc = append(acc, sum)
	}
	fmt.Println(acc)
	return 0
}

// BinarySearch 采用二分查找快速取得抽中的物品的编号
func BinarySearch(arr []float64, target float64) int {
	return 0
}
