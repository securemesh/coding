package codes

type Code struct {
	Value uint32
	Bits  uint8
}

func CodeForIndex(index uint32) Code {
	switch {
	case index < 4:
		return Code{Value: index, Bits: 4}
	case index < 8:
		return Code{Value: 0b01000 + (index & 0b00011), Bits: 5}
	case index < 16:
		return Code{Value: 0b011000 + (index & 0b000111), Bits: 6}
	case index < 32:
		return Code{Value: 0b1000000 + (index & 0b0001111), Bits: 7}
	default:
		set := uint8(index / uint32(64))

		return Code{
			Value: (((2 << set) - 1) << 7) + (index % uint32(64)),
			Bits:  set + 8,
		}
	}
}
