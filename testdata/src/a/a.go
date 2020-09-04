package a

func f1() int {
	num := 10
	ans := 0
	for i := 0; i < num; i++ {
		if i%2 == 0 {
			ans++
		} else if i%3 == 0 {
			ans += 2
		} else {
			ans -= 3
		}
	}
	return ans
}

func f2() int { // want "function f2 is too complicated 12 > 10"
	n := 10
	switch n {
	case 1:
		n++
	case 2:
		n++
	case 3:
		n++
	case 4:
		n++
	case 5:
		n++
	case 6:
		n++
	case 7:
		n++
	case 8:
		n++
	case 9:
		if n%10 == 0 {
			n++
		}
	case 10:
		n++
	default:
		n++
	}
	return n
}

func f3() bool {
	n := 0
	if n%2 == 0{
		n++
	}
	b := n%2 == 0 && n%3 == 0
	return b
}

func f4() int {
	num := []int{1,2,3,4,5}
	ans := 0
	for i := range num {
		if i % 2 == 0 {
			ans += i
		}
	}
	return ans
}
