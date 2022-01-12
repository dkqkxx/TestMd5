package main

import (
	"fmt"
	"math"
)

func Format(n uint) [32]string {
	r := [32]string{}
	for i, c := range fmt.Sprintf("%032b", n) {
		r[31-i] = string(c)
	}
	for i := 0; i < 32; i = i + 8 {
		for j := 0; j < 4; j++ {
			r[i+j], r[i+7-j] = r[i+7-j], r[i+j]
		}
	}
	return r
}

func Bits(bs []byte) []string {
	k := ((len(bs)+8)/64 + 1) * 512
	r := make([]string, 0, k)
	for i := range bs {
		s := fmt.Sprintf("%08b", bs[i])
		for j := range s {
			r = append(r, string(s[j]))
		}
	}
	r = append(r, "1")
	for i := len(r); i < k-64; i++ {
		r = append(r, "0")
	}
	for n := len(bs) * 8; n != 0; n = n / 256 {
		s := fmt.Sprintf("%08b", byte(n%256))
		for j := range s {
			r = append(r, string(s[j]))
		}
	}
	for i := len(r); i < k; i++ {
		r = append(r, "0")
	}
	return r
}

func SUM(a, b [32]string) [32]string {
	r, c := [32]string{}, "0"
	for i := 0; i < len(r)/8; i++ {
		for j := i*8 + 7; j >= i*8; j-- {
			r[j] = XOR(XOR(a[j], b[j]), c)
			c = OR(AND(a[j], b[j]), AND(c, XOR(a[j], b[j])))
		}
	}
	return r
}

func SHIFT(s [32]string, n uint) [32]string {
	nn := int(n % 32)
	for i := 0; i < len(s)/8; i++ {
		for j := 0; j < 8-j; j++ {
			s[i*8+j], s[i*8+7-j] = s[i*8+7-j], s[i*8+j]
		}
	}
	for i := 0; i < len(s)-i; i++ {
		s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
	}
	for i := 0; i < nn-i; i++ {
		s[i], s[nn-1-i] = s[nn-1-i], s[i]
	}
	for i := nn; i < len(s)-i+nn; i++ {
		s[i], s[len(s)-1-i+nn] = s[len(s)-1-i+nn], s[i]
	}
	for i := 0; i < len(s)/8; i++ {
		for j := 0; j < 8-j; j++ {
			s[i*8+j], s[i*8+7-j] = s[i*8+7-j], s[i*8+j]
		}
	}
	return s
}

func F(x, y, z string) string {
	return OR(AND(x, y), AND(NOT(x), z))
}

func G(x, y, z string) string {
	return OR(AND(x, z), AND(y, NOT(z)))
}

func H(x, y, z string) string {
	return XOR(XOR(x, y), z)
}

func I(x, y, z string) string {
	return XOR(y, OR(x, NOT(z)))
}

func LOOP(ABCD [4][32]string, m [16][32]string) [4][32]string {
	f, g, k, s := [32]string{}, 0, [64][32]string{}, [64]uint{}
	for i := 0; i < 64; i++ {
		f[i%32] = "0"
		k[i] = Format(uint(math.Floor(math.Abs(math.Sin(float64(i+1)) * math.Pow(2, 32)))))
		s[i] = [16]uint{7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21}[i/16*4+i%4]
	}
	A, B, C, D := ABCD[0], ABCD[1], ABCD[2], ABCD[3]
	for i := 0; i < 64; i++ {
		if i < 16 {
			for j := 0; j < 32; j++ {
				f[j] = F(B[j], C[j], D[j])
			}
			g = i % 16
		} else if i < 32 {
			for j := 0; j < 32; j++ {
				f[j] = G(B[j], C[j], D[j])
			}
			g = (5*i + 1) % 16
		} else if i < 48 {
			for j := 0; j < 32; j++ {
				f[j] = H(B[j], C[j], D[j])
			}
			g = (3*i + 5) % 16
		} else {
			for j := 0; j < 32; j++ {
				f[j] = I(B[j], C[j], D[j])
			}
			g = (7 * i) % 16
		}
		A, B, C, D = D, SUM(B, SHIFT(SUM(SUM(A, f), SUM(k[i], m[g])), s[i])), B, C
	}
	ABCD[0], ABCD[1], ABCD[2], ABCD[3] = SUM(ABCD[0], A), SUM(ABCD[1], B), SUM(ABCD[2], C), SUM(ABCD[3], D)
	return ABCD
}

func MD5(b []byte) [128]string {
	ABCD, s, r := [4][32]string{Format(0x67452301), Format(0xefcdab89), Format(0x98badcfe), Format(0x10325476)}, Bits(b), [128]string{}
	for i := 0; i < len(s)/512; i++ {
		m := [16][32]string{}
		for j := 0; j < len(m); j++ {
			for k := 0; k < len(m[0]); k++ {
				m[j][k] = s[i*512+j*len(m[0])+k]
			}
		}
		ABCD = LOOP(ABCD, m)
	}
	for i := range ABCD {
		for j := range ABCD[i] {
			r[i*32+j] = ABCD[i][j]
		}
	}
	return r
}
