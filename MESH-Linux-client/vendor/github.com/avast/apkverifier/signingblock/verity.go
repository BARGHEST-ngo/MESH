package signingblock

import (
	"bytes"
	"crypto"
	"fmt"
	"hash"
	"io"
	"math"
)

const (
	verityChunkSize = 4096
	verityAlgo      = crypto.SHA256
)

type verityTreeBuilder struct {
	salt   []byte
	hasher hash.Hash
}

func newVerityTreeBuilder(salt []byte) *verityTreeBuilder {
	return &verityTreeBuilder{
		salt:   salt,
		hasher: verityAlgo.New(),
	}
}

func (b *verityTreeBuilder) generateTreeRootHash(contents dataSource) ([]byte, error) {
	digestSize := int64(b.hasher.Size())
	levelOffset, err := b.calculateLevelOffset(contents.length(), digestSize)
	if err != nil {
		return nil, err
	}

	verityBuffer := make([]byte, levelOffset[len(levelOffset)-1])
	for i := len(levelOffset) - 2; i >= 0; i-- {
		sink := byteWriter{dest: verityBuffer[levelOffset[i]:levelOffset[i+1]]}

		var src dataSource
		if i == len(levelOffset)-2 {
			src = contents
			if err := b.digestDataByChunks(src, &sink); err != nil {
				return nil, err
			}
		} else {
			src = &dataSourceBytes{verityBuffer[levelOffset[i+1]:levelOffset[i+2]]}
			if err := b.digestDataByChunks(src, &sink); err != nil {
				return nil, err
			}
		}

		totalOutput := b.divideRoundup(src.length(), verityChunkSize) * digestSize
		incomplete := totalOutput % verityChunkSize
		if incomplete > 0 {
			padding := make([]byte, verityChunkSize-incomplete)
			if _, err := sink.Write(padding); err != nil {
				return nil, err
			}
		}
	}

	firstPage := verityBuffer[:verityChunkSize]
	return b.saltedDigest(firstPage), nil
}

func (b *verityTreeBuilder) calculateLevelOffset(dataSize int64, digestSize int64) ([]int32, error) {
	var levelSize []int64
	for {
		chunkCount := b.divideRoundup(dataSize, verityChunkSize)
		size := verityChunkSize * b.divideRoundup(chunkCount*digestSize, verityChunkSize)
		levelSize = append(levelSize, size)
		if chunkCount*digestSize <= verityChunkSize {
			break
		}
		dataSize = chunkCount * digestSize
	}

	levelOffset := make([]int32, len(levelSize)+1)
	for i := 0; i < len(levelSize); i++ {
		size := levelSize[len(levelSize)-i-1]
		if size < math.MinInt32 || size > math.MaxInt32 {
			return nil, fmt.Errorf("verity level size overflow")
		}
		levelOffset[i+1] = levelOffset[i] + int32(size)
	}
	return levelOffset, nil
}

func (b *verityTreeBuilder) digestDataByChunks(src dataSource, sink io.Writer) error {
	size := src.length()
	var offset int64
	buf := bytes.NewBuffer(make([]byte, 0, verityChunkSize))
	for ; offset+verityChunkSize <= size; offset += verityChunkSize {
		buf.Reset()
		if err := src.writeTo(buf, offset, verityChunkSize); err != nil {
			return err
		}

		hash := b.saltedDigest(buf.Bytes())
		if _, err := sink.Write(hash); err != nil {
			return err
		}
	}

	remaining := size % verityChunkSize
	if remaining > 0 {
		buf.Reset()
		if err := src.writeTo(buf, offset, remaining); err != nil {
			return err
		}

		// padding
		buf.Write(make([]byte, verityChunkSize-remaining))

		hash := b.saltedDigest(buf.Bytes())
		if _, err := sink.Write(hash); err != nil {
			return err
		}
	}
	return nil
}

func (b *verityTreeBuilder) saltedDigest(data []byte) []byte {
	b.hasher.Reset()
	if len(b.salt) != 0 {
		b.hasher.Write(b.salt)
	}
	b.hasher.Write(data)
	return b.hasher.Sum(nil)
}

func (b *verityTreeBuilder) divideRoundup(dividend, divisor int64) int64 {
	return (dividend + divisor - 1) / divisor
}
