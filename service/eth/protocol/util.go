package protocol

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type hexInt int
type hexBig big.Int

func (b *hexBig) AsBigInt() *big.Int {
	return (*big.Int)(b)
}

func (b *hexInt) UnmarshalJSON(input []byte) error {
	str := string(input)
	str = strings.Trim(str, `"`)
	str = strings.TrimPrefix(str, "0x")

	dec, err := strconv.ParseUint(str, 16, 64)
	*b = hexInt(dec)
	return err
}

func (b *hexBig) UnmarshalJSON(input []byte) error {
	str := string(input)
	str = strings.Trim(str, `"`)

	i := big.Int{}
	_, err := fmt.Sscan(str, &i)
	*b = hexBig(i)
	return err
}

// ParseInt parse hex string value to int
func ParseInt(value string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(value, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// ParseBigInt parse hex string value to big.Int
func ParseBigInt(value string) (big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(value, &i)

	return i, err
}

// IntToHex convert int to hexadecimal representation
func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}

// BigToHex covert big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}
	return "0x" + strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0")
}
