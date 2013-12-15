package utils

// utility class that contains bits
type BitList struct {
	count int
	data  []int32
}

// returns a new BitList with the given length
// all bits are initialize with false
func NewBitList(capacity int) *BitList {
	bl := new(BitList)
	bl.count = capacity
	x := 0
	if capacity%32 != 0 {
		x = 1
	}
	bl.data = make([]int32, capacity/32+x)
	return bl
}

// returns the number of contained bits
func (bl *BitList) Len() int {
	return bl.count
}

func (bl *BitList) grow() {
	growBy := len(bl.data)
	if growBy < 128 {
		growBy = 128
	} else if growBy >= 1024 {
		growBy = 1024
	}

	nd := make([]int32, len(bl.data)+growBy)
	copy(nd, bl.data)
	bl.data = nd
}

// appends the given bit to the end of the list
func (bl *BitList) AddBit(bit bool) {
	itmIndex := bl.count / 32
	for itmIndex >= len(bl.data) {
		bl.grow()
	}
	bl.SetBit(bl.count, bit)
	bl.count++
}

// sets the bit at the given index to the given value
func (bl *BitList) SetBit(index int, value bool) {
	itmIndex := index / 32
	itmBitShift := 31 - (index % 32)
	if value {
		bl.data[itmIndex] = bl.data[itmIndex] | 1<<uint(itmBitShift)
	} else {
		bl.data[itmIndex] = bl.data[itmIndex] & ^(1 << uint(itmBitShift))
	}
}

// returns the bit at the given index
func (bl *BitList) GetBit(index int) bool {
	itmIndex := index / 32
	itmBitShift := 31 - (index % 32)
	return ((bl.data[itmIndex] >> uint(itmBitShift)) & 1) == 1
}

// appends all 8 bits of the given byte to the end of the list
func (bl *BitList) AddByte(b byte) {
	for i := 7; i >= 0; i-- {
		bl.AddBit(((b >> uint(i)) & 1) == 1)
	}
}

// appends the last (LSB) 'count' bits of 'b' the the end of the list
func (bl *BitList) AddBits(b int, count byte) {
	for i := int(count - 1); i >= 0; i-- {
		bl.AddBit(((b >> uint(i)) & 1) == 1)
	}
}

// appends all bits in the bool-slice to the end of the list
func (bl *BitList) AddRange(bits []bool) {
	for _, b := range bits {
		bl.AddBit(b)
	}
}

// returns all bits of the BitList as a []byte
func (bl *BitList) GetBytes() []byte {
	len := bl.count >> 3
	if (bl.count % 8) != 0 {
		len += 1
	}
	result := make([]byte, len)
	for i := 0; i < len; i++ {
		shift := (3 - (i % 4)) * 8
		result[i] = (byte)((bl.data[i/4] >> uint(shift)) & 0xFF)
	}
	return result
}

// iterates through all bytes contained in the BitList
func (bl *BitList) IterateBytes() <-chan byte {
	res := make(chan byte)

	go func() {
		c := bl.count
		shift := 24
		i := 0
		for c > 0 {
			res <- byte((bl.data[i] >> uint(shift)) & 0xFF)
			shift -= 8
			if shift < 0 {
				shift = 24
				i += 1
			}
			c -= 8
		}
		close(res)
	}()

	return res
}
