package lottery

import (
	"fmt"
	"testing"
)

func TestLottery(t *testing.T) {
	items := []int{2, 2, 4, 7}
	fmt.Println(Lottery(items))
}

func TestBinarySearch(t *testing.T) {
	items := []int{2, 4, 6, 9, 15}
	fmt.Println(BinarySearch(items, -1))
	fmt.Println(BinarySearch(items, 2))
	fmt.Println(BinarySearch(items, 4))
	fmt.Println(BinarySearch(items, 6))
	fmt.Println(BinarySearch(items, 9))
	fmt.Println(BinarySearch(items, 15))
	fmt.Println(BinarySearch(items, 18))
}
