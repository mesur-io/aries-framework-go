/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package bbs12381g2pub

import (
	"encoding/binary"
	"errors"
)

func uint32ToBytes(value uint32) []byte {
	bytes := make([]byte, 4)

	binary.BigEndian.PutUint32(bytes, value)

	return bytes
}

func uint16FromBytes(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func uint32FromBytes(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func bitvectorToIndexes(data []byte) []int {
	revealedIndexes := make([]int, 0)
	scalar := 0

	for _, v := range data {
		remaining := 8

		for v > 0 {
			revealed := v & 1
			if revealed == 1 {
				revealedIndexes = append(revealedIndexes, scalar)
			}

			v >>= 1
			scalar++
			remaining--
		}

		scalar += remaining
	}

	return revealedIndexes
}

type pokPayload struct {
	messagesCount int
	revealed      []int
}

// nolint:gomnd
func parsePoKPayload(bytes []byte) (*pokPayload, error) {
	if len(bytes) < 2 {
		return nil, errors.New("invalid size of PoK payload")
	}

	messagesCount := int(uint16FromBytes(bytes[0:2]))
	offset := lenInBytes(messagesCount)
	revealed := bitvectorToIndexes(reverseBytes(bytes[2:offset]))

	return &pokPayload{
		messagesCount: messagesCount,
		revealed:      revealed,
	}, nil
}

// nolint:gomnd
func (p *pokPayload) toBytes() []byte {
	bytes := make([]byte, p.lenInBytes())

	binary.BigEndian.PutUint16(bytes, uint16(p.messagesCount))

	bitvector := bytes[2:]

	for _, r := range p.revealed {
		idx := r / 8
		bit := r % 8

		bitvector[idx] |= 1 << bit
	}

	reverseBytes(bitvector)

	return bytes
}

func (p *pokPayload) lenInBytes() int {
	return lenInBytes(p.messagesCount)
}

func lenInBytes(messagesCount int) int {
	return 2 + (messagesCount / 8) + 1
}

func newPoKPayload(messagesCount int, revealed []int) *pokPayload {
	return &pokPayload{
		messagesCount: messagesCount,
		revealed:      revealed,
	}
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
