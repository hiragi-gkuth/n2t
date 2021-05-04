package main

import (
	"io/ioutil"
	"strings"
)

type CommandType uint8

const (
	A_COMMAND = 0b01
	C_COMMAND = 0b10
	L_COMMAND = 0b11
)

type IParser interface {
	HasMoreCommands() bool
	Advance()
	CommandType() CommandType
	Symbol() string
	Dest() string
	Comp() string
	Jump() string
	ResetSeeker()
}

type Parser struct {
	commands []string
	seeker   int
}

func NewParser(asmPath string) IParser {
	rawAsmCode, e := ioutil.ReadFile(asmPath)
	if e != nil {
		panic(e)
	}
	commands := asmRawCodeToCommands(rawAsmCode)

	return &Parser{
		commands: commands,
		seeker:   0,
	}
}

func (p Parser) HasMoreCommands() bool {
	return p.seeker != (len(p.commands) - 1)
}

func (p *Parser) Advance() {
	p.seeker++
}

func (p *Parser) ResetSeeker() {
	p.seeker = 0
}

func (p Parser) CommandType() CommandType {
	nowCommand := p.commands[p.seeker]
	if strings.HasPrefix(nowCommand, "@") {
		return A_COMMAND
	}
	if strings.HasPrefix(nowCommand, "(") {
		return L_COMMAND
	}
	return C_COMMAND
}

func (p Parser) Symbol() string {
	nowCommand := p.commands[p.seeker]
	if p.CommandType() == A_COMMAND {
		symbol := strings.TrimPrefix(nowCommand, "@")
		return symbol
	}
	if p.CommandType() == L_COMMAND {
		symbol := strings.TrimPrefix(nowCommand, "(")
		symbol = strings.TrimSuffix(symbol, ")")
		return symbol
	}

	// cannot be reach here
	panic("Symbol() cannot be called when command type is C")
}

func (p Parser) Dest() string {
	if p.CommandType() != C_COMMAND {
		panic("Dest() cannot be called when command type is not C")
	}

	nowCommand := p.commands[p.seeker]
	equalOperatorPos := strings.Index(nowCommand, "=")
	if equalOperatorPos == -1 {
		return "_"
	}
	return nowCommand[0:equalOperatorPos]
}

func (p Parser) Comp() string {
	if p.CommandType() != C_COMMAND {
		panic("Comp() cannot be called when command type is not C")
	}

	nowCommand := p.commands[p.seeker]
	equalOperatorPos := strings.Index(nowCommand, "=")
	semicolonOperatorPos := strings.Index(nowCommand, ";")

	// return as it is when Dest and Jump is omitted
	if equalOperatorPos == -1 && semicolonOperatorPos == -1 {
		return nowCommand
	}
	if equalOperatorPos == -1 && semicolonOperatorPos != -1 {
		return nowCommand[:semicolonOperatorPos]
	}
	if equalOperatorPos != -1 && semicolonOperatorPos == -1 {
		return nowCommand[equalOperatorPos+1:]
	}
	return nowCommand[equalOperatorPos+1 : semicolonOperatorPos]
}

func (p Parser) Jump() string {
	if p.CommandType() != C_COMMAND {
		panic("Jump() cannot be called when command type is not C")
	}

	nowCommand := p.commands[p.seeker]
	semicolonOperatorPos := strings.Index(nowCommand, ";")
	if semicolonOperatorPos == -1 {
		return ""
	}
	return nowCommand[semicolonOperatorPos+1:]
}

func asmRawCodeToCommands(rawCode []byte) (commands []string) {
	rawCodeArray := strings.Split(string(rawCode), "\n")

	// skipping comments and empty line
	for _, line := range rawCodeArray {
		line = strings.TrimSpace(line)

		// skip comment
		commentPos := strings.Index(line, "//")
		if commentPos != -1 {
			line = line[:commentPos]
		}

		// skip empty line
		if strings.Compare(line, "") == 0 {
			continue
		}
		commands = append(commands, line)
	}
	return
}
