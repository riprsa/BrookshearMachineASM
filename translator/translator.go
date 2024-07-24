package assembly

import (
	"encoding/hex"
	"errors"
	"log"
	"log/slog"
	"strings"
)

// adds instructions as comments
var DEBUG = true

type Program struct {
	Words    [256]byte
	Last     byte
	Comments [128]string

	// vars stores address of variable
	Vars map[string]byte
	// labels stores address of label
	Labels map[string]byte
}

func (p *Program) editByte(src byte, value byte) {
	p.Words[int(src)] = value
}

func (p Program) Bytes() []byte {
	var s string

	for i := 0; i < len(p.Words); i += 2 {
		s += hex.EncodeToString([]byte{p.Words[i], p.Words[i+1]})
		if p.Comments[i/2] != "" {
			s += ";" + p.Comments[i/2]
		}
		s += "\n"
	}

	return []byte(s)
}

type CleanLine struct {
	line    string
	comment string
}

func NewProgram(fileString string) Program {
	p := Program{
		Vars:   map[string]byte{},
		Labels: map[string]byte{},
	}

	rawLines := strings.Split(fileString, "\n")

	var cleanLines []CleanLine

	// to fill labels
	for _, rawLine := range rawLines {

		// remove whitespaces for labels (and not only)
		trimmed := strings.TrimSpace(rawLine)

		// remove empty and comma lines
		if trimmed == "" {
			continue
		}

		if string(trimmed[0]) == ";" {
			// TODO: save separated comments
			continue
		}

		// remove multiple spaces
		trimmed = strings.Join(strings.Fields(trimmed), " ")

		// split a word by space
		bySemicolon := strings.Split(trimmed, ";")

		var commentLine string
		if len(bySemicolon) > 1 {
			trimmed = bySemicolon[0]
			// trim again to remove spaces before comments
			trimmed = strings.TrimSpace(trimmed)

			// join bySemicolon without first element
			bySemicolon = bySemicolon[1:]
			commentLine = strings.Join(bySemicolon, ";")
			commentLine = strings.TrimSpace(commentLine)
		}

		cl := CleanLine{
			line:    trimmed,
			comment: commentLine,
		}

		cleanLines = append(cleanLines, cl)
	}

	var beforeLabelVars []string
	var last byte
	// handle labels
	for _, cl := range cleanLines {
		cmd := strings.Split(cl.line, " ")
		if len(cmd) == 0 {
			slog.Error("empty line")
			continue
		}

		switch cmd[0] {
		case ".label":
			name := strings.TrimSuffix(cmd[1], ":")
			p.Labels[name] = last
			if DEBUG {
				if cl.comment != "" {
					p.Comments[last/2] = "// " + cl.comment + " "
				}
				p.Comments[last/2] = ".label " + name + ": "
			} else {
				if cl.comment != "" {
					p.Comments[last/2] = "// " + cl.comment + " "
				}
			}

		case "var":
			if len(cmd) == 4 { // var A F0 0x0F
				name := cmd[1]
				address := cmd[2]

				addressBytes, err := hex.DecodeString(address)
				if err != nil {
					panic(err)
				}
				addressByte := addressBytes[0]

				value := cmd[3]
				v, had := strings.CutPrefix(value, "0x")
				if !had {
					panic("values shoud be in format 0xFF: " + cmd[3])
				}
				valueBytes, err := hex.DecodeString(v)
				if err != nil {
					panic(err)
				}
				valueByte := valueBytes[0]

				p.editByte(addressByte, valueByte)
				p.Vars[name] = addressByte

				if DEBUG {
					p.Comments[addressByte/2] = cl.line
				}

				continue
			} else if len(cmd) == 3 { // var A LABEL_NAME
				beforeLabelVars = append(beforeLabelVars, cl.line)
				continue
			}
			log.Fatal("var should have 2 or 3 arguments")
		case "array":
			if len(cmd) < 3 {
				log.Fatal("array should have 2 arguments, example: array A0 0x0F")
			}
			address := cmd[1]
			addressBytes, err := hex.DecodeString(address)
			if err != nil {
				panic(err)
			}
			addressByte := addressBytes[0]

			if DEBUG {
				p.Comments[addressByte/2] = cl.line
			}

			var values []byte
			for _, value := range cmd[2:] {
				v, had := strings.CutPrefix(value, "0x")
				if !had {
					panic("array: values shoud be in format 0xFF: " + cmd[3])
				}
				valueBytes, err := hex.DecodeString(v)
				if err != nil {
					panic(err)
				}
				values = append(values, valueBytes[0])
			}

			for i, value := range values {
				p.editByte(addressByte+byte(i), value)
			}
		default:
			beforeLabelVars = append(beforeLabelVars, cl.line)
			if DEBUG {
				if cl.comment != "" {
					cl.comment = cl.line + " // " + cl.comment
				} else {
					cl.comment = cl.line
				}
			} else {
				if cl.comment != "" {
					cl.comment = "// " + cl.comment
				}
			}
			p.Comments[last/2] = p.Comments[last/2] + cl.comment

			last += 2
		}
	}

	var afterVarLabels []string
	for _, line := range beforeLabelVars {

		cmd := strings.Split(line, " ")
		if len(cmd) == 0 {
			slog.Error("empty line")
			continue
		}

		if cmd[0] == "var" {
			if len(cmd) != 3 {
				log.Fatal("var should have 2 arguments, example: var A LABEL_NAME")
			}

			name := cmd[1]
			labelName := cmd[2]

			addressByte, ok := p.Labels[labelName]
			if !ok {

				panic("label is missing: " + labelName)
			}

			p.Vars[name] = addressByte + 1

			continue
		}

		afterVarLabels = append(afterVarLabels, line)
	}

	p.Last = 0
	for _, line := range afterVarLabels {
		// remove whitespaces for labels (and not only)
		line := strings.TrimSpace(line)

		// split a word by space
		cmd := strings.Split(line, " ")

		// by first arg:
		switch cmd[0] {
		case "load":
			addressByte, ok := p.Vars[cmd[1]]
			if !ok {
				address := cmd[1]
				addressBytes, err := hex.DecodeString(address)
				if err != nil {
					panic(err)
				}
				addressByte = addressBytes[0]
			}

			r, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			dest, err := hex.DecodeString("1" + r)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, dest[0]) // to concatinate
			p.editByte(p.Last+1, addressByte)

			p.Last += 2
		case "loadi":
			f, had := strings.CutPrefix(cmd[1], "0x")
			if !had {
				panic("values shoud be in format 0xFF: " + cmd[3])
			}

			value, err := hex.DecodeString(f)
			if err != nil {
				panic(err)
			}

			r, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			dest, err := hex.DecodeString("2" + r)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, dest[0]) // to concatinate
			p.editByte(p.Last+1, value[0])

			p.Last += 2

		case "store":
			addressByte, ok := p.Vars[cmd[2]]
			if !ok {
				address := cmd[1]
				addressBytes, err := hex.DecodeString(address)
				if err != nil {
					panic(err)
				}
				addressByte = addressBytes[0]
			}

			r, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			dest, err := hex.DecodeString("3" + r)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, dest[0])
			p.editByte(p.Last+1, addressByte)

			p.Last += 2

		case "mov":
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("40" + r1 + r2)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "add":
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			r3, err := register(cmd[3])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("5" + r3 + r2 + r1)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "addf": // NOTE: may be not working
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			r3, err := register(cmd[3])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("6" + r3 + r2 + r1)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "or":
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			r3, err := register(cmd[3])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("7" + r3 + r2 + r1)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "and":
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			r3, err := register(cmd[3])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("8" + r3 + r2 + r1)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "xor":
			r1, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			r2, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			r3, err := register(cmd[3])
			if err != nil {
				panic(err)
			}

			bs, err := hex.DecodeString("9" + r3 + r2 + r1)
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "rotate":
			r, err := register(cmd[1])
			if err != nil {
				panic(err)
			}

			// TODO: make number 0x0 format instead of 0
			bs, err := hex.DecodeString("A" + r + "0" + cmd[2])
			if err != nil {
				panic(err)
			}

			p.editByte(p.Last, bs[0])
			p.editByte(p.Last+1, bs[1])

			p.Last += 2

		case "jeq":
			r, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			value, err := hex.DecodeString("B" + r)
			if err != nil {
				panic(err)
			}

			addressByte, ok := p.Labels[cmd[1]]
			if !ok {
				address := cmd[1]
				addressBytes, err := hex.DecodeString(address)
				if err != nil {
					panic(err)
				}
				addressByte = addressBytes[0]
			}

			p.editByte(p.Last, value[0])
			p.editByte(p.Last+1, addressByte)

			p.Last += 2

		case "halt":

			p.editByte(p.Last, 0xC0)
			p.Last += 2
		case "jgt":
			r, err := register(cmd[2])
			if err != nil {
				panic(err)
			}

			value, err := hex.DecodeString("D" + r)
			if err != nil {
				panic(err)
			}

			addressByte, ok := p.Labels[cmd[1]]
			if !ok {
				address := cmd[1]
				addressBytes, err := hex.DecodeString(address)
				if err != nil {
					panic(err)
				}
				addressByte = addressBytes[0]
			}

			p.editByte(p.Last, value[0])
			p.editByte(p.Last+1, addressByte)

			p.Last += 2
		}
	}
	return p
}

func register(r string) (string, error) {
	a, ok := strings.CutPrefix(r, "r")
	if !ok {
		return "", errors.New("registers start with 'r': " + r)
	}
	return a, nil
}
