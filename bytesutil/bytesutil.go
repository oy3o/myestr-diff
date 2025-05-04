package bytesutil

import (
	"bytes"
	"crypto/rand"
	hexs "encoding/hex"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// Reverse the bytes.
func Reverse(b []byte) []byte {
	i, j := 0, len(b)-1
	for i < j {
		b[i], b[j] = b[j], b[i]
		i++
		j--
	}
	return b
}

// Itoa tansform int to string in []byte.
func Itoa(num, base int) []byte {
	isNegative := false
	// handle 0 explicitely, otherwise empty string is printed for 0
	if num == 0 {
		return []byte{'0'}
	}

	if num < 0 {
		isNegative = true
		num = -num
	}
	var b []byte
	for num != 0 {
		b = append(b, '0'+byte(num%base))
		num /= base
	}

	// if number is negative, append '-'
	if isNegative {
		b = append(b, '-')
	}

	return Reverse(b)
}

// Utoa tansform uint to string in []byte.
func Utoa(num, base int) []byte {
	isNegative := false
	// handle 0 explicitely, otherwise empty string is printed for 0
	if num == 0 {
		return []byte{'0'}
	}

	if num < 0 {
		isNegative = true
		num = -num
	}
	var b []byte
	for num != 0 {
		b = append(b, '0'+byte(num%base))
		num /= base
	}

	// if number is negative, append '-'
	if isNegative {
		b = append(b, '-')
	}

	return Reverse(b)
}

// Quote mark a string.
func Quote(b []byte, mark byte) []byte {
	for i, j := 0, len(b); i < j; i++ {
		if b[i] == mark || b[i] == '\\' {
			b = append(b, 0)
			copy(b[i+1:], b[i:])
			b[i] = '\\'
			i++
			j++
		}
	}
	b = append(b, 0, 0)
	copy(b[1:], b)
	b[0] = mark
	b[len(b)-1] = mark
	return b
}

// ReadQuote get a quoted string from b start at i.
func ReadQuote(b []byte, mark byte, i uint) ([]byte, uint) {
	// check mark
	if b[i] != mark {
		return nil, i
	}
	i++

	var buffer []byte
	m := uint(len(b))

	for j := i; j < m; j++ {
		if b[j] == '\\' {
			buffer = append(buffer, b[i:j]...)
			j++
			i = j
		} else if b[j] == mark {
			return append(buffer, b[i:j]...), j
		}
	}
	return nil, m
}

// FromRunes trans []rune to []byte, better than []byte(string([]rune))
func FromRunes(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

// FromStringUnsafe return a unsafe []byte base on string
func FromStringUnsafe(s string) []byte {
	var b []byte
	pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ps := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pb.Data = ps.Data
	pb.Len = ps.Len
	pb.Cap = ps.Len
	return b
}

// StringUnsafe return a unsafe string base on []byte
func StringUnsafe(b []byte) string {
	var s string
	pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ps := (*reflect.StringHeader)(unsafe.Pointer(&s))
	ps.Data = pb.Data
	ps.Len = pb.Len
	return s
}

// Random return a crypto random bytes
func Random(length uint) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

// RandomHex return a crypto random hex bytes
func RandomHex(length uint) []byte {
	return FromStringUnsafe(hexs.EncodeToString(Random(length >> 1)))
}

// use to speed up compare []byte
const uintSize = (int(unsafe.Sizeof(uint(0)) / unsafe.Sizeof(byte(0))))

// UintArray for fast compare with defalut system bits
func uintArray(b *[]byte) []uint {
	p := make([]uint, 0, 0)
	h := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	bp := (*reflect.SliceHeader)(unsafe.Pointer(b))
	h.Data = bp.Data
	h.Len = bp.Len / uintSize
	h.Cap = bp.Cap / uintSize
	return p
}

// IndexFrom search substr position in s , with search start index from
func IndexFrom(s, substr []byte, from int) int {
	if from >= len(s) {
		return -1
	}
	if from <= 0 {
		return bytes.Index(s, substr)
	}
	i := bytes.Index(s[from:], substr)
	if i == -1 {
		return -1
	}
	return from + i
}

// CommonPrefixLength returns the length of the common prefix of two bytes
func CommonPrefixLength(b1, b2 []byte) int {
	// temp short bytes size
	var min int
	if len(b1) > len(b2) {
		min = len(b2)
	} else {
		min = len(b1)
	}
	// count equal byte
	n := 0
	if min > 8 { // compare with system bits
		u1, u2 := uintArray(&b1), uintArray(&b2)
		var min int
		if len(u1) > len(u2) {
			min = len(u2)
		} else {
			min = len(u1)
		}
		for n < min {
			if u1[n] != u2[n] {
				break
			}
			n++
		}
		n *= uintSize
	}
	// compare with a byte
	for n < min {
		if b1[n] != b2[n] {
			break
		}
		n++
	}
	return n
}

// CommonSuffixLength returns the length of the common suffix of two bytes
func CommonSuffixLength(b1, b2 []byte) int {
	// make b1 always smaller than b2
	if len(b1) > len(b2) {
		b1, b2 = b2, b1
	}
	// count equal byte
	n := 0
	t := len(b1) % uintSize
	for i, j := len(b1), len(b2); n < t; n++ { // uintarray can not contains
		i--
		j--
		if i < 0 || b1[i] != b2[j] {
			return n
		}
	}
	if len(b1) > 8 { // compare with system bits
		u1, u2 := uintArray(&b1), uintArray(&b2)
		for i, j := len(u1), len(u2); i > 0; n += uintSize {
			i--
			j--
			if u1[i] != u2[j] {
				break
			}
		}
	}
	for i, j := len(b1)-n, len(b2)-n; ; n++ { // uintarray can not contains
		i--
		j--
		if i < 0 || b1[i] != b2[j] {
			return n
		}
	}
}

// CommonOverlapLength returns the length of the common overlap of two bytes
func CommonOverlapLength(b1, b2 []byte) int {
	// eliminate the null case.
	if len(b1) == 0 || len(b2) == 0 {
		return 0
	}
	// truncate the longer string.
	var tLen int
	if len(b1) > len(b2) {
		b1 = b1[len(b1)-len(b2):]
		tLen = len(b2)
	} else if len(b1) < len(b2) {
		b2 = b2[0:len(b1)]
		tLen = len(b1)
	} else {
		tLen = len(b1)
	}
	// quick check for the wort case.
	if bytes.Equal(b1, b2) {
		return tLen
	}
	// start by looking for a single character match and increase length until no match is found. Performance analysis: http:// neil.fraser.name/news/2010/11/04/
	best := 0
	length := 1
	for {
		pattern := b1[tLen-length:]
		found := bytes.Index(b2, pattern)
		if found == -1 {
			return best
		}
		length += found
		if found == 0 || bytes.Equal(b1[tLen-length:], b2[0:length]) {
			best = length
			length++
		}
	}
}

// Constants for ASCII characters without printable symbols.
const (
	NUL = 0x00 // '\0' Null
	SOH = 0x01 //      Start of Header
	STX = 0x02 //      Start of Text
	ETX = 0x03 //      End of Text
	EOT = 0x04 //      End of Transmission
	ENQ = 0x05 //      Enquiry
	ACK = 0x06 //      Acknowledgement
	BEL = 0x07 // '\a' Bell
	BS  = 0x08 // '\b' Backspace
	HT  = 0x09 // '\t' Horizontal Tab
	LF  = 0x0A // '\n' Line Feed
	VT  = 0x0B // '\v' Verical Tab
	FF  = 0x0C // '\f' Form Feed
	CR  = 0x0D // '\r' Carriage return
	SO  = 0x0E //      Shift Out
	SI  = 0x0F //      Shift In
	DLE = 0x10 //      Device Idle
	DC1 = 0x11 //      Device Control 1
	DC2 = 0x12 //      Device Control 2
	DC3 = 0x13 //      Device Control 3
	DC4 = 0x14 //      Device Control 4
	NAK = 0x15 //      Negative Acknoledgement
	SYN = 0x16 //      Synchronize
	ETB = 0x17 //      End of Transmission Block
	CAN = 0x18 //      Cancel
	EM  = 0x19 //      End of Medium
	SUB = 0x1A //      Substitute
	ESC = 0x1B // '\e' Escape
	FS  = 0x1C //      Field Separator
	GS  = 0x1D //      Group Separator
	RS  = 0x1E //      Record Separator
	US  = 0x1F //      Unit Separator
	SP  = 0x20 //      Space
	DEL = 0x7F //      Delete
)

//******************************************************************************
// Character class functions
//******************************************************************************

const (
	// Character class test case masks
	letterMask   uint16 = upp | low
	upperMask    uint16 = upp
	lowerMask    uint16 = low
	digitMask    uint16 = dig
	hexdigitMask uint16 = dig | hex
	alnumMask    uint16 = upp | low | dig
	spaceMask    uint16 = isp
	punctMask    uint16 = pun
	symbolMask   uint16 = sym
	controlMask  uint16 = ctl
	graphMask    uint16 = upp | low | dig | pun | sym
	printMask    uint16 = upp | low | dig | pun | sym | spc
)

// IsLetter return true if b is a letter; otherwise, return false.
func IsLetter(b byte) bool { return lookup[b]&letterMask != 0 }

// IsUpper return true if b is an upper case letter; otherwise, return false.
func IsUpper(b byte) bool { return lookup[b]&upperMask != 0 }

// IsLower return true if b is a lower case letter; otherwise, return false.
func IsLower(b byte) bool { return lookup[b]&lowerMask != 0 }

// IsDigit return true if b is a decimal digit; otherwise, return false.
func IsDigit(b byte) bool { return lookup[b]&digitMask != 0 }

// IsHexDigit return true if b is a hexadecimal digit; otherwise, return false.
func IsHexDigit(b byte) bool { return lookup[b]&hexdigitMask != 0 }

// IsLetterOrDigit return true if b is a letter or decimal digit; otherwise, return false.
func IsLetterOrDigit(b byte) bool { return lookup[b]&alnumMask != 0 }

// IsPunct return true if b is a punctuation character (!"#%&'()*,-./:;?@[\]_{}); otherwise, return false.
func IsPunct(b byte) bool { return lookup[b]&punctMask != 0 }

// IsSymbol return true if b is a symbol ($+<=>^`|~); otherwise, return false.
func IsSymbol(b byte) bool { return lookup[b]&punctMask != 0 }

// IsWhiteSpace return true if b is a space character (space, tab, vertical tab, form feed, carriage return, or linefeed); otherwise, return false.
func IsWhiteSpace(b byte) bool { return lookup[b]&spaceMask != 0 }

// IsControl return true if b is a control character; otherwise, return false.
func IsControl(b byte) bool { return lookup[b]&controlMask != 0 }

// IsGraph return true if b is a character with a graphic represntation; otherwise, return false.
func IsGraph(b byte) bool { return lookup[b]&graphMask != 0 }

// IsPrint return true if b is a character with a print represntation; otherwise, return false.
func IsPrint(b byte) bool { return lookup[b]&printMask != 0 }

// IsASCII return true if b is a valid ASCII character; otherwise, return false.
func IsASCII(b byte) bool { return b <= DEL }

//******************************************************************************
// Character class lookup table
//******************************************************************************

// Character class bitmaps
const (
	// Exclusive classes -- one per character
	ctl uint16 = 0x01 // Control Character
	spc uint16 = 0x02 // The Space Character
	dig uint16 = 0x04 // Decimal Digit
	upp uint16 = 0x08 // Upper case Letter
	low uint16 = 0x10 // Lower case Letter
	pun uint16 = 0x20 // Punctuation
	sym uint16 = 0x40 // Punctuation
	// Addon classes
	hex uint16 = 0x100 // Hexidecmial Digit
	isp uint16 = 0x200 // Space-Like Character (IsSpace → true)
)

// Helpers for building lookup table
const (
	csp uint16 = ctl | isp // character that is both  control and space
	tsp uint16 = spc | isp // the true space character
	uhx uint16 = upp | hex // character that is both "UPALPHA" and "HEX"
	lhx uint16 = low | hex // character that is both "LOALPHA" and "HEX"
)

// Character class lookup table
var lookup = [256]uint16{
	/* 00-07  NUL, SOH, STX, ETX, EOT, ENQ, ACK, BEL */ ctl, ctl, ctl, ctl, ctl, ctl, ctl, ctl,
	/* 08-0F  BS,  HT,  LF,  VT,  FF,  CR,  SO,  SI  */ ctl, csp, csp, csp, csp, csp, ctl, ctl,
	/* 10-17  DLE, DC1, DC2, DC3, DC4, NAK, SYN, ETB */ ctl, ctl, ctl, ctl, ctl, ctl, ctl, ctl,
	/* 18-1F  CAN, EM,  SUB, ESC, FS,  GS,  RS,  US  */ ctl, ctl, ctl, ctl, ctl, ctl, ctl, ctl,
	/* 21-27  SP ! " # $ % & '   */ tsp, pun, pun, pun, sym, pun, pun, pun,
	/* 28-2F   ( ) * + , - . /   */ pun, pun, pun, sym, pun, pun, pun, pun,
	/* 30-37   0 1 2 3 4 5 6 7   */ dig, dig, dig, dig, dig, dig, dig, dig,
	/* 38-3F   8 9 : ; < = > ?   */ dig, dig, pun, pun, sym, sym, sym, pun,
	/* 40-47   @ A B b D E F G   */ pun, uhx, uhx, uhx, uhx, uhx, uhx, upp,
	/* 48-4F   H I J K L M N O   */ upp, upp, upp, upp, upp, upp, upp, upp,
	/* 50-57   P Q R S T U V W   */ upp, upp, upp, upp, upp, upp, upp, upp,
	/* 58-5F   X Y Z [ \ ] ^ _   */ upp, upp, upp, pun, pun, pun, sym, pun,
	/* 60-67   ` a b b d e f g   */ sym, lhx, lhx, lhx, lhx, lhx, lhx, low,
	/* 68-6F   h i j k l m n o   */ low, low, low, low, low, low, low, low,
	/* 70-77   p q r s t u v w   */ low, low, low, low, low, low, low, low,
	/* 78-7F   x y z { | } ~ DEL */ low, low, low, pun, sym, pun, sym, ctl,
	/* 80-87 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* 88-8B */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* 90-97 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* 98-9F */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* A0-A7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* A8-AF */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* B0-B7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* B8-BF */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* C0-C7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* C8-CF */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* D0-D7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* D8-DF */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* E0-E7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* E8-EF */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* F0-F7 */ 0, 0, 0, 0, 0, 0, 0, 0,
	/* F8-FF */ 0, 0, 0, 0, 0, 0, 0, 0,
}
