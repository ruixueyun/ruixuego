// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package bytepool

import (
	"fmt"
)

const (
	MaxIndex = 31 // NewMultiRatedBytePool 内各长度池的数量, 区间为 [1, 32]
	Magic    = 0x077CB531
)

// RatedBytePool 定宽字节对象池
// 一般用在有对字节切换进行反复读写, 且数据长度固定的场景下, 比如读取协议头
type RatedBytePool struct {
	c chan []byte
	w int
}

// NewRatedBytePool 新建一个定宽字节对象池
func NewRatedBytePool(maxSize, width int) (bp *RatedBytePool) {
	return &RatedBytePool{
		c: make(chan []byte, maxSize),
		w: width,
	}
}

// Get 从池中获取元素
func (rbp *RatedBytePool) Get() (b []byte) {
	select {
	case b = <-rbp.c:
	default:
		b = make([]byte, rbp.w)
	}
	return
}

// Put 将元素放回池中
func (rbp *RatedBytePool) Put(b []byte) {
	if cap(b) < rbp.w {
		return
	}

	select {
	case rbp.c <- b[:rbp.w]:
	default:
	}
}

// Capacity 获取当前池的容量
func (rbp *RatedBytePool) Capacity() int {
	return len(rbp.c)
}

// Width 获取当前池内元素的宽度
func (rbp *RatedBytePool) Width() (n int) {
	return rbp.w
}

// Reset 重制元素数据
func (*RatedBytePool) Reset(b []byte) []byte {
	return b[:0]
}

// MultiRatedBytePool
// 该池为定宽字节对象池的一个扩展实现, 可根据要存取数据的长获得一个足够容量的定额元素.
type MultiRatedBytePool struct {
	p      []*RatedBytePool
	min    int64
	max    int64
	offset int
}

// NewMultiRatedBytePool 创建扩展定宽字节对象池
// 注: 池中每个元素的大小都是 2^N
// param: minSize 池中最小元素的大小
// param: count RatedBytePool 的数量
// param: capacityPerItemPool 每个 RatedBytePool 的元素数量
func NewMultiRatedBytePool(
	minSize, count, capacityPerItemPool int) (*MultiRatedBytePool, error) {
	ret := &MultiRatedBytePool{}

	ret.min = Pow2roundup(int64(minSize))
	ret.offset = log2(ret.min)       // 算出最小元素其在 debruijn 的索引数
	if ret.offset+count > MaxIndex { // 如果最大索引超过阈值(debruijn的上限)则报错
		return nil, fmt.Errorf(
			"when minSize is %d, count cannot be greater than %d",
			minSize, MaxIndex-ret.offset)
	}

	ret.p = make([]*RatedBytePool, count)
	startSize := int(ret.min)
	for i := 0; i < count; i++ {
		ret.p[i] = NewRatedBytePool(capacityPerItemPool, startSize)
		startSize <<= 1 // 每个 NewRatedBytePool 元素宽度是前一个的 2 倍
	}
	ret.max = int64(startSize >> 1)
	return ret, nil
}

// MinCapacity 获取池中最小规格元素的容量, 小于这个容量将直接取用最小容量的元素
func (mrbp *MultiRatedBytePool) MinCapacity() int64 {
	return mrbp.min
}

// MaxCapacity 获取池中最大规格元素的容量, 超过这个容量无法使用池
func (mrbp *MultiRatedBytePool) MaxCapacity() int64 {
	return mrbp.max
}

// Get 获取一个足够存取 size 长度的元素, 最终的元素长度 >= size
func (mrbp *MultiRatedBytePool) Get(size int) []byte {
	capacity := Pow2roundup(int64(size))
	if capacity > mrbp.max {
		// 超过池中元素的最大容量时直接新建, 此处不能利用池的特性
		return make([]byte, size, capacity)
	}
	if capacity < mrbp.min {
		// 小于最小容量则直接拿一个最小容量值的
		capacity = mrbp.min
	}
	// log2(capacity)-mrbp.offset 是为了让其从池的第 0 号下标开始,
	// 否则因为最小池为 minSize 字节, 下标 0 ~ mrbp.offset 永远无法被使用
	return mrbp.p[log2(capacity)-mrbp.offset].Get()
}

// Put 将元素放回池中
func (mrbp *MultiRatedBytePool) Put(bytes []byte) bool {
	capacity := int64(cap(bytes))
	if capacity < mrbp.min || capacity > mrbp.max {
		return false
	}
	mrbp.p[log2(capacity)-mrbp.offset].Put(bytes)
	return true
}

// 德布莱英序列
var debruijn = [...]int{
	0, 1, 28, 2, 29, 14, 24, 3,
	30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7,
	26, 12, 18, 6, 11, 5, 10, 9,
}

// log2 计算一个数的二进制中从低位到高位第一个 1 经过了多少个 0 (末尾 0 的个数)
func log2(x int64) int {
	return debruijn[uint32(x*Magic)>>27]
}

// Pow2roundup 获取大于等于某个数最近的一个 2^N
func Pow2roundup(x int64) int64 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	return x + 1
}
