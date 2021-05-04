package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	parser := NewParser(os.Args[1])

	// pass 1
	symbolTable := symbolLabelResolver(parser)

	// pass 2
	hackBinaryCode := assemble(parser, symbolTable)

	for _, b := range hackBinaryCode {
		fmt.Printf("%016b\n", b)
	}
}

func symbolLabelResolver(parser IParser) ISymbolTable {
	symbolTable := NewSymbolTable()
	romAddr := uint16(0x0)
	for {
		cType := parser.CommandType()

		if cType == A_COMMAND || cType == C_COMMAND {
			romAddr++
		}
		if cType == L_COMMAND {
			symbolTable.AddLabelEntry(parser.Symbol(), romAddr)
		}

		if !parser.HasMoreCommands() {
			break
		}
		parser.Advance()
	}
	return symbolTable
}

func assemble(parser IParser, symbolTable ISymbolTable) []uint16 {
	parser.ResetSeeker()
	coder := NewCode()
	hackBinaryCode := []uint16{}

	for {
		cType := parser.CommandType()
		// do not assemble L_COMMAND on pass 2
		if cType == L_COMMAND {
			if !parser.HasMoreCommands() {
				break
			}
			parser.Advance()
			continue
		}

		var commandBits uint16
		switch {
		case cType == C_COMMAND:
			commandBits = assembleCCommand(
				coder.Dest(parser.Dest()),
				coder.Jump(parser.Jump()),
				coder.Comp(parser.Comp()),
			)
		case cType == A_COMMAND:
			commandBits = assembleACommand(symbolTable, parser.Symbol())
		default:
			panic("unrecognized commands")
		}

		hackBinaryCode = append(hackBinaryCode, commandBits)
		if !parser.HasMoreCommands() {
			break
		}
		parser.Advance()
	}
	return hackBinaryCode
}

func assembleCCommand(dest, jump, comp uint16) (commandBits uint16) {
	commandBitsBase := uint16(0b1110_0000_0000_0000)
	commandBits = commandBitsBase | jump | (dest << 3) | (comp << 6)
	return
}

func assembleACommand(symbolTable ISymbolTable, symbol string) (commandBits uint16) {
	symbolValue, e := strconv.Atoi(symbol)
	if e == nil {
		// use symbol value as commandBits if symbol is immediate value
		commandBits = uint16(symbolValue)
		return
	}

	// symbol resolver
	if !symbolTable.Contains(symbol) {
		symbolTable.AddVariableEntry(symbol)
	}
	commandBits = symbolTable.GetAddress(symbol)
	return
}
