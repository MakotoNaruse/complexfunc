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

func f2() int { // want "function a.f2 is too complicated 12 > 10"
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

func f3(n int) bool {
	if n%2 == 0{
		return true
	}
	if n % 3 == 0 {
		return true
	}
	return false
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

func f5() int {
	ans := 0
	var ch chan int
	select {
	case v := <- ch:
		ans = v
	default:
		ans = -1
	}
	return ans
}

func f6() { // want "function a.f6 has redundant branch"
	a := 0
	if false {
	}
	a++
}

func f7() int {
	n := 10
	var a int
	switch {
	case 0 < n && n < 100, n > 200:
		a = 1
	}
	return a
}
