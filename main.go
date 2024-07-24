package main

import (
	"bmasm/runtime"
	assembly "bmasm/translator"
	"flag"
	"fmt"
	"os"
)

func main() { // go run main.go -asm=test.asm -o=bytecode.txt -e
	// Define command-line flags
	bytecodeFlag := flag.String("bytecode", "", "Path to the bytecode file")
	asmFlag := flag.String("asm", "", "Path to the assembly file")

	eFlag := flag.Bool("e", false, "Do excute asm after and do print RAM compactly")
	outputFlag := flag.String("o", "", "Name of the output file")

	// Parse command-line flags
	flag.Parse()

	// Check if either bytecode or assembly file is provided
	if *bytecodeFlag == "" && *asmFlag == "" {
		fmt.Println("Please provide either a bytecode or assembly file")
		return
	}

	// Check if both bytecode and assembly file are provided
	if *bytecodeFlag != "" && *asmFlag != "" {
		fmt.Println("Please provide either a bytecode or assembly file, not both")
		return
	}

	if *bytecodeFlag != "" {
		m, err := runtime.NewFromFile(*bytecodeFlag)
		if err != nil {
			panic(err)
		}

		err = m.Run()
		if err != nil {
			panic(err)
		}

		if *outputFlag != "" {
			os.WriteFile(*outputFlag, []byte(m.Output()), 0655)
		}
		return
	}

	if *asmFlag != "" {
		content, err := os.ReadFile(*asmFlag)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		p := assembly.NewProgram(string(content))

		if *eFlag {
			m, err := runtime.New(string(p.Bytes()))
			if err != nil {
				panic(err)
			}

			fmt.Println(
				"initial:",
			)
			for k, v := range p.Vars {
				b := m.ByAddress(v)
				fmt.Printf("  %s: %s => %s => 0b%s => 0x%s\n", k, b.Decimal(), b.Decimal2s(), b.Binary(), b.String())
			}

			os.WriteFile("initial.txt", p.Bytes(), 0655)

			err = m.Run()
			if err != nil {
				panic(err)
			}

			fmt.Println(
				"final:",
			)
			for k, v := range p.Vars {
				b := m.ByAddress(v)
				fmt.Printf("  %s: %s => %s => 0b%s => 0x%s\n", k, b.Decimal(), b.Decimal2s(), b.Binary(), b.String())
			}

			fmt.Print(m.String())

			if *outputFlag != "" {
				os.WriteFile(*outputFlag, []byte(m.Output()), 0655)
			}

		} else {
			// Write the buffer content to the file
			err = os.WriteFile(*outputFlag, p.Bytes(), 0655)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
}
