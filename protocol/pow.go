package protocol

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/PoW-HC/hashcash/pkg/pow"
)

var (
	errWrongFormat = errors.New("wrong msg format")
)

// Message - string presentation of hashcash
// Format  - 1:bits:date:resource:externsion:rand:counter
func Unmarshal(message string) (*pow.Hashcach, error) {

	attrs := strings.Split(message, ":")
	if len(attrs) != 7 {
		return nil, errWrongFormat
	}
	// challenge.Version


	ver, err := strconv.Atoi(attrs[0])
	if err != nil {
		return nil, err
	}
	bits, err := strconv.Atoi(attrs[1])
	if err != nil {
		return nil, err
	}
	dateStr, err := strconv.ParseInt(attrs[2], 10, 64)
	if err != nil {
		return nil, err
	}
	date := time.Unix(dateStr, 0)
	resource := attrs[3]
	ext := attrs[4]
	rand, err := base64.StdEncoding.DecodeString(attrs[5])
	if err != nil {
		return nil, err
	}
	counterStr, err := base64.StdEncoding.DecodeString(attrs[6])
	if err != nil {
		return nil, err
	}
	counter, err := strconv.ParseInt(string(counterStr), 16, 64)
	if err != nil {
		return nil, err
	}
	hash := pow.NewHashcach(int32(ver), int32(bits), date, resource, ext, rand, counter)
	return hash, nil
}
