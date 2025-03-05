//go:build !solution

package speller

var (
	figures      = []string{"zero", "one ", "two ", "three ", "four ", "five ", "six ", "seven ", "eight ", "nine "}
	lessTwenties = []string{"ten ", "eleven ", "twelve ", "thirteen ", "fourteen ", "fifteen ", "sixteen ", "seventeen ", "eighteen ", "nineteen "}
	tens         = []string{"ten ", "twenty ", "thirty ", "forty ", "fifty ", "sixty ", "seventy ", "eighty ", "ninety "}
	large        = []string{"hundred ", "thousand ", "million ", "billion "}
)

func make3(n int64, res string) (int64, string) {
	for n > 0 {
		if n < 10 {
			res += figures[n]
			n = 0
			return n, res
		} else if n < 20 {
			res += lessTwenties[n-10]
			n = 0
			return n, res
		} else if n < 100 {
			res += tens[n/10-1]
			n = n % 10
			if n != 0 {
				n, res = make3(n, res[:len(res)-1]+"-")
			} else {
				n, res = make3(n, res)
			}
		} else {
			_, res = make3(n/100, res)
			res += large[0]
			n, res = make3(n%100, res)
		}
	}
	return n, res
}

func makeAll(n int64) string {
	res, globRes, minus, i := "", "", "", 0
	if n < 0 {
		n *= -1
		minus = "minus "
	}
	for n > 0 {
		_, res = make3(n%1000, "")
		if i > 0 && res != "" {
			globRes = large[i] + globRes
		}
		globRes = res + globRes
		n /= 1000
		i++
	}
	return minus + globRes
}

func Spell(n int64) string {
	if n == 0 {
		return figures[0]
	}
	res := makeAll(n)
	return res[:len(res)-1]
}
