package protocol

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PoW-HC/hashcash/pkg/pow"
)

type PoW struct{}


// Message - string presentation of hashcash
// Format  - 1:bits:date:resource:externsion:rand:counter
func (_ PoW) Unmarshal(message string) *pow.Hashcach {
	attrs := strings.Split(message, ":")
	fmt.Println(attrs)
	if len(attrs) != 7 {
		return nil
	}
	ver, _ := strconv.Atoi(attrs[0])
	bits, _ := strconv.Atoi(attrs[1])
	dateStr, _ := strconv.ParseInt(attrs[2], 10, 64)
	date := time.Unix(dateStr, 0)
	resource := attrs[3]
	ext := attrs[4]
	rand, _ := base64.StdEncoding.DecodeString(attrs[5])
	counterStr, _ := base64.StdEncoding.DecodeString(attrs[6])
	counter, _ := strconv.ParseInt(string(counterStr), 16, 64)

	hash := pow.NewHashcach(int32(ver), int32(bits), date, resource, ext, rand, counter)
	return hash
}
