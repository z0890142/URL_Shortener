package common

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encode(num uint32) string {
	encoded := ""
	for num > 0 {
		remainder := num % 62
		num /= 62
		encoded = string(base62Chars[remainder]) + encoded
	}
	return encoded
}
