package bitwise

// Set a flag
func Set(n uint32, flag Flag) uint32 { return n | uint32(flag) }

// Clear a flag
func Clear(n uint32, flag Flag) uint32 { return n &^ uint32(flag) }

// Toggle a flag
func Toggle(n uint32, flag Flag) uint32 { return n ^ uint32(flag) }

// Has a flag
func Has(n uint32, flag Flag) bool { return n&uint32(flag) == uint32(flag) }

// Swap values so that a = b and b = a
func Swap(a uint32, b uint32) (uint32, uint32) {
	a ^= b
	b ^= a
	a ^= b
	return a, b
}
