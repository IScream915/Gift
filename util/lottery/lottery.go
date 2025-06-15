package lottery

// Lottery 抽奖函数
// 入参: 抽奖物品的库存量
// 出参: 抽中的物品的编号
func Lottery(probs []float64) int {
	if len(probs) == 0 {
		return -1 // 没有参与抽奖的物品, 返回-1报错
	}
	sum := 0.0
	return 0
}
