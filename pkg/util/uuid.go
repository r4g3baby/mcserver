package util

import (
	"crypto/md5"
	"github.com/google/uuid"
)

func NameUUIDFromBytes(data []byte) uuid.UUID {
	uniqueId := md5.Sum(data)
	uniqueId[6] = (uniqueId[6] & 0x0f) | uint8((3&0xf)<<4) // version 3
	uniqueId[8] = (uniqueId[8] & 0x3f) | 0x80              // RFC 4122 variant
	return uniqueId
}
