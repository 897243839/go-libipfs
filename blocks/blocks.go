// Package blocks contains the lowest level of IPLD data structures.
// A block is raw data accompanied by a CID. The CID contains the multihash
// corresponding to the block.
package blocks

import (
	"errors"
	"fmt"
	hc "github.com/897243839/HcdComp"
	cid "github.com/ipfs/go-cid"
	dshelp "github.com/ipfs/go-ipfs-ds-help"
	u "github.com/ipfs/go-ipfs-util"
	mh "github.com/multiformats/go-multihash"
	"time"
)

// ErrWrongHash is returned when the Cid of a block is not the expected
// according to the contents. It is currently used only when debugging.
var ErrWrongHash = errors.New("data did not match given hash")

// Block provides abstraction for blocks implementations.
type Block interface {
	RawData() []byte
	Cid() cid.Cid
	String() string
	Loggable() map[string]interface{}
}

// A BasicBlock is a singular block of data in ipfs. It implements the Block
// interface.
type BasicBlock struct {
	cid  cid.Cid
	data []byte
}

// NewBlock creates a Block object from opaque data. It will hash the data.
func NewBlock(data []byte) *BasicBlock {
	// TODO: fix assumptions

	println("NewBlock-压缩类型", hc.GetCompressorType(data))
	return &BasicBlock{data: data, cid: cid.NewCidV0(u.Hash(data))}
}
func AutoHC(data []byte, c cid.Cid) []byte {
	//start := time.Now().UnixNano()
	key := dshelp.MultihashToDsKey(c.Hash()).String()[1:]
	startTime3 := time.Now().UnixNano()
	tp := hc.GetCompressorType(data)
	endTime3 := time.Now().UnixNano()
	dur3 := endTime3 - startTime3
	println("gettype-time:", dur3)
	//println(tp)
	if tp != hc.UnknownCompressor {
		println("1-ZLIB", tp)
		startTime1 := time.Now().UnixNano()
		data = hc.Decompress(data, tp)
		v, ok := hc.MapLit.Get(key)
		if !ok {
			hc.MapLit.Set(key, 1)
		}
		if v == 5 {
			hc.Tsf <- key
			v += 1
			hc.MapLit.Set(key, v)
		} else if v < 5 {
			v += 1
			hc.MapLit.Set(key, v)
		}
		endTime1 := time.Now().UnixNano()
		dur1 := endTime1 - startTime1
		println("mapcool-time:", dur1)
		return data
	} else {
		startTime := time.Now().UnixNano()
		v, ok := hc.Maphot.Get(key)
		if !ok {
			hc.Maphot.Set(key, 1)
		}
		if v > 999 {
		} else {
			v += 1
			hc.Maphot.Set(key, v)
		}
		endTime := time.Now().UnixNano()
		dur := endTime - startTime
		println("maphot-time:", dur)
		println("0-ipfs", tp)
	}
	//endT := time.Now().UnixNano()
	//durt := endT - start
	//println("sum-time:", durt)
	//println("0-ipfs", tp)
	return data
}

// NewBlockWithCid creates a new block when the hash of the data
// is already known, this is used to save time in situations where
// we are able to be confident that the data is correct.
func NewBlockWithCid(data []byte, c cid.Cid) (*BasicBlock, error) {
	//println("NewBlockWithCid.cid=", dshelp.MultihashToDsKey(c.Hash()).String()[1:])
	if u.Debug {
		chkc, err := c.Prefix().Sum(data)
		if err != nil {
			return nil, err
		}

		if !chkc.Equals(c) {
			return nil, ErrWrongHash
		}
	}
	data = AutoHC(data, c)
	return &BasicBlock{data: data, cid: c}, nil
}

func NewBlockWithCid1(data []byte, c cid.Cid) (*BasicBlock, error) {
	//println("NewBlockWithCid.cid=", dshelp.MultihashToDsKey(c.Hash()).String()[1:])
	if u.Debug {
		chkc, err := c.Prefix().Sum(data)
		if err != nil {
			return nil, err
		}
		if !chkc.Equals(c) {
			return nil, ErrWrongHash
		}
	}
	return &BasicBlock{data: data, cid: c}, nil
}

// Multihash returns the hash contained in the block CID.
func (b *BasicBlock) Multihash() mh.Multihash {
	return b.cid.Hash()
}

// RawData returns the block raw contents as a byte slice.
func (b *BasicBlock) RawData() []byte {
	return b.data
}

// Cid returns the content identifier of the block.
func (b *BasicBlock) Cid() cid.Cid {
	return b.cid
}

// String provides a human-readable representation of the block CID.
func (b *BasicBlock) String() string {
	return fmt.Sprintf("[Block %s]", b.Cid())
}

// Loggable returns a go-log loggable item.
func (b *BasicBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"block": b.Cid().String(),
	}
}
