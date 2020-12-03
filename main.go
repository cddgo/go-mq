package gomq
//
//import (
//	"fmt"
//)
//
//func main() {
//	topics := make(map[string][]chan int)
//
//	ch1 := make(chan int, 1)
//	ch2 := make(chan int, 2)
//	ch3 := make(chan int, 3)
//	ch4 := make(chan int, 4)
//	//ch5 := make(chan int, 4)
//	chssss := []chan int{ch1, ch2, ch3, ch4}
//	topics["sub1"] = chssss
//
//	chs := topics["sub1"]
//	fmt.Printf("%p\n", chs)
//	for i, ch := range chs {
//		if ch == ch4 {
//			copy(chs[i:], chs[i+1:])
//			chs = chs[:len(chs)-1]
//			break
//		}
//	}
//
//	newchs := topics["sub1"]
//	fmt.Printf("%p\n", newchs)
//	for i := 0; i < 3; i++ {
//		fmt.Println(chs[i] == newchs[i])
//	}
//
//	for _, ch := range newchs {
//		fmt.Println(cap(ch))
//	}
//}
