package hotp

import (
	"URL_Shortener/pkg/utils/algorithm"
	"crypto/hmac"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type GenerateOptions struct {
	Algorithm algorithm.Algorithm
	Digits    int
}

func HotpGenerateCode(secret string, counter uint64, opts GenerateOptions) (string, error) {
	secret = strings.ToUpper(secret)

	// secretBytes, err := base32.StdEncoding.DecodeString(secret)
	// if err != nil {
	// 	return "", fmt.Errorf("HotpGenerateCode: %s", err.Error())
	// }

	buf := make([]byte, 8)
	mac := hmac.New(opts.Algorithm.Hash, []byte(secret))
	binary.BigEndian.PutUint64(buf, counter)

	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits
	mod := int32(value % int64(math.Pow10(l)))
	f := fmt.Sprintf("%%0%dd", l)
	return fmt.Sprintf(f, mod), nil
}
