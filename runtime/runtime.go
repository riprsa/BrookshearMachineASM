package runtime

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

type Value byte

func (v Value) String() string {
	return fmt.Sprintf("%X", byte(v))
}

func (v Value) Decimal2s() string {
	return fmt.Sprintf("%d", int8(v))
}

func (v Value) Decimal() string {
	return fmt.Sprintf("%d", v)
}

func (v Value) Binary() string {
	return fmt.Sprintf("%08b", v)
}

type Machine struct {
	Memory [256]Value
	Reg    [16]Value

	// next index to execute
	Counter Value
}

func NewFromFile(path string) (*Machine, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(string(bytes))
}

func New(data string) (*Machine, error) {
	final := [256]Value{}

	rawLines := strings.Split(data, "\n")

	c := 0
	for _, raw := range rawLines {
		if raw == "" {
			continue
		}

		ss := strings.Split(raw, ";")
		i := ss[0]
		if len(i) != 4 {
			return nil, fmt.Errorf("bad bytes: %s", raw)
		}

		bytes, err := hex.DecodeString(i)
		if err != nil {
			return nil, err
		}
		if len(bytes) != 2 {
			return nil, fmt.Errorf("bad bytes: %s", bytes)
		}

		final[c] = Value(bytes[0])
		final[c+1] = Value(bytes[1])

		c += 2
	}

	return &Machine{
		Memory: final,
	}, nil
}

func (m Machine) ByAddress(a byte) Value {
	return m.Memory[a]
}

func (m Machine) String() string {
	ram := m.Memory

	lines := " | 0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n"
	for i := 0; i <= 15; i++ {
		line := string(hex.EncodeToString([]byte{byte(i) << 4})[0]) + "| "
		for j := 0; j <= 15; j++ {
			current := 16*i + j

			line += hex.EncodeToString([]byte{byte(ram[current])}) + " "
		}
		lines += line + "\n"
	}

	return lines
}

func (m Machine) Output() string {
	ram := m.Memory

	lines := ""
	for i := 0; i <= 255; i += 2 {
		lines += hex.EncodeToString([]byte{
			byte(ram[i]), byte(ram[i+1]),
		}) + "\n"
	}

	return lines
}

func (m *Machine) Run() error {
loop:
	for {

		b1, b2 := m.Memory[m.Counter], m.Memory[m.Counter+1]
		instr := newInstruction(b1, b2)

		// fmt.Println("counter: ", m.Counter)
		// fmt.Println("instruction: ", instr)

		switch instr.opcode {
		case 1: // memory -> register
			m.Reg[instr.half1] = m.Memory[b2]
		case 2: // value -> register
			m.Reg[instr.half1] = b2
		case 3: // reg -> memory
			m.Memory[b2] = m.Reg[instr.half1]
		case 4: // reg -> reg
			m.Reg[instr.half3] = m.Reg[instr.half2]
		case 5: // add 2s
			m.Reg[instr.half1] = m.Reg[instr.half2] + m.Reg[instr.half3]
		case 6: // add float // TODO:
			log.Fatal("i cant make it work..")
			m.Reg[instr.half1] = m.Reg[instr.half2] + m.Reg[instr.half3]
		case 7:
			m.Reg[instr.half1] = m.Reg[instr.half2] | m.Reg[instr.half3]
		case 8:
			m.Reg[instr.half1] = m.Reg[instr.half2] & m.Reg[instr.half3]
		case 9:
			m.Reg[instr.half1] = m.Reg[instr.half2] ^ m.Reg[instr.half3]
		case 10: // "Rotate bits in register %X cyclically right %X steps", i.half1, i.half3
			b := m.Reg[instr.half1]
			n := instr.half3

			for range n {
				rotated := b >> 1
				rotated |= b << 7
				b = rotated
			}

			m.Reg[instr.half1] = b

		case 11: // "Jump to cell %X%X if register %X equals register 0", i.half2, i.half3, i.half1
			if m.Reg[instr.half1] == m.Reg[0] {
				m.Counter = b2
				continue
			}
		case 12:
			break loop
		case 13: // D RXY
			// Jump to instruction in RAM cell XY if the content of register R is greater than (>) the content of register 0. Data is interpreted as integers in two's-complement notation.

			if m.Reg[instr.half1] > m.Reg[0] {
				m.Counter = m.Memory[b2]
			}

		default:
			return fmt.Errorf("wrong opcode: %d", instr.opcode)
		}

		m.Counter += 2
		// time.Sleep(time.Millisecond * 100)
	}

	return nil
}

func (m *Machine) Reset() {
	m = &Machine{}
}

type Instruction struct {
	opcode int
	half1  int
	half2  int
	half3  int
}

func newInstruction(fb, sb Value) Instruction {
	return Instruction{
		opcode: int(fb >> 4),
		half1:  int(fb & 0x0F),
		half2:  int(sb >> 4),
		half3:  int(sb & 0x0F),
	}
}

func (i Instruction) String() string {
	switch i.opcode {
	case 1:
		return fmt.Sprintf("Copy bits at cell %X%X to register %X", i.half2, i.half3, i.half1)
	case 2:
		return fmt.Sprintf("Copy bit-string %X%X to register %X", i.half2, i.half3, i.half1)
	case 3:
		return fmt.Sprintf("Copy bits in register %X to cell %X%X", i.half1, i.half2, i.half3)
	case 4:
		return fmt.Sprintf("Copy bits in register %X to register %X", i.half2, i.half3)
	case 5:
		return fmt.Sprintf("Add bits in registers %X and %X (two's-complement), put in register %X", i.half2, i.half3, i.half1)
	case 6:
		return fmt.Sprintf("Add bits in registers %X and %X (float), put in register %X", i.half2, i.half3, i.half1)
	case 7:
		return fmt.Sprintf("Bitwise OR bits in registers %X and %X, put in register %X", i.half2, i.half3, i.half1)
	case 8:
		return fmt.Sprintf("Bitwise AND bits in registers %X and %X, put in register %X", i.half2, i.half3, i.half1)
	case 9:
		return fmt.Sprintf("Bitwise XOR bits in register %X and %X, put in register %X", i.half2, i.half3, i.half1)
	case 10:
		return fmt.Sprintf("Rotate bits in register %X cyclically right %X steps", i.half1, i.half3)
	case 11:
		return fmt.Sprintf("Jump to cell %X%X if register %X equals register 0", i.half2, i.half3, i.half1)
	case 12:
		return "Halt"
	case 13:
		return fmt.Sprintf("Jump to cell %X%X if register %X is greater than register 0", i.half2, i.half3, i.half1)
	}

	return "Error"
}

// func main() {
// 	m, err := New("./input.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	m.Run("./output.txt")

// }
