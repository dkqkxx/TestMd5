package main

import (
	"sort"
	"strings"
)

func Level(exp string) int {
	if strings.ContainsAny(exp, "+") {
		return 2
	}
	if strings.ContainsAny(exp, "*") {
		return 1
	}
	return 0
}

func AND(exp1, exp2 string) string {
	res := ""
	switch Level(exp1)*3 + Level(exp2) {
	case 0: // 00
		res = and00(exp1, exp2)
	case 1: // 01
		res = and01(exp1, exp2)
	case 2: // 02
		res = and02(exp1, exp2)
	case 4: // 11
		res = and11(exp1, exp2)
	case 5: // 12
		res = and12(exp1, exp2)
	case 8: // 22
		res = and22(exp1, exp2)
	default: // 10 20 21
		res = AND(exp2, exp1)
	}
	return res
}

func NOT(exp string) string {
	res := ""
	switch Level(exp) {
	case 0: // 0
		res = not0(exp)
	case 1: // 1
		res = not1(exp)
	case 2: // 2
		res = not2(exp)
	}
	return res
}

func OR(exp1, exp2 string) string {
	return NOT(AND(NOT(exp1), NOT(exp2)))
}

func XOR(exp1, exp2 string) string {
	return OR(AND(NOT(exp1), exp2), AND(exp1, NOT(exp2)))
}

func and00(s1, s2 string) string {
	if s1 == "0" || s2 == "0" {
		return "0"
	}
	if s1 == "1" {
		return s2
	}
	if s2 == "1" {
		return s1
	}
	if s1 == s2 {
		return s1
	}
	if s1 == not0(s2) {
		return "0"
	}
	if s1 < s2 {
		return s1 + "*" + s2
	}
	return s2 + "*" + s1
}

func and01(s1, s2 string) string {
	if s1 == "0" {
		return "0"
	}
	if s1 == "1" {
		return s2
	}
	cs := strings.Split(s2, "*")
	ns1 := not0(s1)
	for i := range cs {
		if cs[i] == s1 {
			return s2
		}
		if cs[i] == ns1 {
			return "0"
		}
	}
	cs = append(cs, s1)
	sort.Strings(cs)
	return strings.Join(cs, "*")
}

func and02(s1, s2 string) string {
	if s1 == "0" {
		return "0"
	}
	if s1 == "1" {
		return s2
	}
	// 原始 s2=A+B不能合并,s1*(A+B)也不能
	is := strings.Split(s2, "+")
	res := ""
	for i := 0; i < len(is); i++ {
		ai := and01(s1, is[i])
		if ai == "0" {
			continue
		}
		if ai == "1" {
			return "1"
		}
		res = res + "+" + ai
	}
	if len(res) == 0 {
		return "0"
	}
	res = res[1:]
	is = strings.Split(res, "+")
	sort.Strings(is)
	return strings.Join(is, "+")
}

func and11(s1, s2 string) string {
	if s1 == s2 {
		return s1
	}
	cs := strings.Split(s1+"*"+s2, "*")
	sort.Strings(cs)
	res := cs[0]
	for i := 1; i < len(cs); i++ {
		if cs[i] == cs[i-1] {
			continue
		}
		if cs[i] == not0(cs[i-1]) {
			return "0"
		}
		res = res + "*" + cs[i]
	}
	return res
}

func and12(s1, s2 string) string {
	is := strings.Split(s2, "+")
	ies := ""
	for i := range is {
		ci := and11(s1, is[i])
		if ci == "0" {
			continue
		}
		if ci == "1" {
			return "1"
		}
		ies = ies + "+" + ci
	}
	if len(ies) == 0 {
		return "0"
	}
	is = strings.Split(ies[1:], "+")
	sort.Strings(is)
	ls := make([]bool, len(is))
	for i := 1; i < len(is); i++ {
		for j := 0; j < i; j++ {
			ai := and11(is[i], is[j])
			if ai == is[i] {
				ls[i] = true
			} else if ai == is[j] {
				ls[j] = true
			}
		}
	}
	res := ""
	for i := 0; i < len(is); i++ {
		if !ls[i] {
			res = res + "+" + is[i]
		}
	}
	res = res[1:] //len(res)>0
	rs := strings.Split(res, "+")
	sort.Strings(rs)
	return strings.Join(rs, "+")
}

func and22(s1, s2 string) string {
	if s1 == s2 {
		return s1
	}
	cs1 := strings.Split(s1, "+")
	cs2 := strings.Split(s2, "+")
	ies := ""
	for i := 0; i < len(cs1); i++ {
		for j := 0; j < len(cs2); j++ {
			is := and11(cs1[i], cs2[j])
			if is == "0" {
				continue
			}
			if is == "1" {
				return "1"
			}
			ies = ies + "+" + is
		}
	}
	if len(ies) == 0 {
		return "0"
	}
	cs := strings.Split(ies[1:], "+")
	ls := make([]bool, len(cs))
	for i := 1; i < len(cs); i++ {
		for j := 0; j < i; j++ {
			ai := and11(cs[i], cs[j])
			if ai == cs[i] {
				ls[i] = true
			} else if ai == cs[j] {
				ls[j] = true
			}
		}
	}
	res := ""
	for i := 0; i < len(cs); i++ {
		if !ls[i] {
			res = res + "+" + cs[i]
		}
	}
	res = res[1:] //len(res)>0
	rs := strings.Split(res, "+")
	sort.Strings(rs)
	return strings.Join(rs, "+")
}

func not0(s string) string {
	if s == "0" {
		return "1"
	}
	if s == "1" {
		return "0"
	}
	if s[len(s)-1:] == "!" {
		return s[:len(s)-1]
	}
	return s + "!"
}

func not1(s string) string {
	cs := strings.Split(s, "*")
	for i := range cs {
		cs[i] = not0(cs[i]) // 不会出现0,1
	}
	return strings.Join(cs, "+")
}

func not2(s string) string {
	is := strings.Split(s, "+")
	res := not1(is[0])
	for i := range is {
		res = and22(res, not1(is[i]))
	}
	return res
}
