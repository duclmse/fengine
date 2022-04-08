package main

import (
	"fmt"
	"strconv"
)

// 1234235 => 8
// 1/ 1 2 3 4 235
// 2/ 1 23 4235
// 3/ 1 234 235
// 4/ 1 234235
// 5/ 12 34 235
// 6/ 12 34235
// 7/ 123 4235
// 8/ 1234235

var length int
var out int

func main() {
	num := "9999999999999"
	fmt.Printf("num='%s'\n", num)
	fmt.Printf("count=%d; out=%d\n", numberOfCombinations(num), out)
}

func numberOfCombinations(num string) int {
	out = 0
	length = len(num)
	count := 0
	loop := 0
	for digit := 1; digit <= length; digit++ {
		if digit > length {
			fmt.Printf("length=%d; digit=%d", length, digit)
			continue
		}
		loop += 1
		valid := float64(length) / float64(digit)
		if valid < 2 {
			break
		}
		count += split(num, digit, 0, 0, 0, "")
		fmt.Printf(">>> loop=%d\n", loop)
	}
	return out
}

func split(num string, digit int, minix int, from int, level int, output string) int {
	//fmt.Printf("digit=%d, minix=%5d, from=%2d; s=%s\n", digit, minix, from, string(num[from]))
	d := digit
	if length < from+d-1 || num[from] == '0' {
		return 0
	}
	prev := minix
	count := 0
	for {
		i := from
		current := selectOne(num, &d, &i, prev, level)
		fmt.Printf(" => %5d (i=%d)", current, i)
		if current == 0 {
			return 0
		}
		if i+digit < length {
			fmt.Printf(" => split(num=%7s, digit=%d, minix=%d, i=%d, lv=%d)\n", num[i:], digit, current, i, level+1)
			count += split(num, digit, current, i, level+1, fmt.Sprintf("%s %d", output, current))
			//fmt.Printf(" split count=%d; lv=%d\n", count, level)
			d++
		} else {
			fmt.Printf("\n:> %s %d\n", output, current)
			out += 1
			return count
		}
	}
}

func selectOne(num string, digit *int, i *int, prev int, level int) int {
	fmt.Printf(" => s1 num=%8s digit=%d; i=%d; lv=%d\n", num[*i:], *digit, *i, level)
	to := *digit + *i
	if to > length {
		return 0
	}
	for {
		current, _ := convert(num, *i, to)
		if prev > current {
			*digit += 1
			to += 1
			continue
		}
		if float64(length-*i-*digit)/float64(*digit) < 2 {
			remain, _ := convert(num, *i, length)
			*i += length - *i
			return remain
		}
		*i += *digit
		return current
	}
}

func convert(num string, from int, to int) (int, error) {
	return strconv.Atoi(num[from:to])
}
