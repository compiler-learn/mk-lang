package code

// package code implements the bytecode instruction set for our virtual
// machine that will execute monkey-lan source code.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Definition struct {
	Name          string
	OperandWidths []int
}

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

func (o Opcode) String() string {
	def, err := Lookup(byte(o))
	if err != nil {
		return ""
	}
	return def.Name
}

const (
	// 加载常量
	LoadConstant Opcode = iota
	// 加载内置
	LoadBuiltin
	// AssignGlobal ...
	// 全局分配
	AssignGlobal
	// AssignLocal ...
	// 局部分配
	AssignLocal
	// 加载全局
	LoadGlobal
	// 绑定全局
	BindGlobal
	// 加载局部
	LoadLocal
	BindLocal
	// 加载free
	LoadFree
	// 加载模块
	LoadModule
	// ?
	SetSelf
	// true
	LoadTrue
	LoadFalse
	LoadNull
	// 获取元素
	GetItem
	// 设置元素
	SetItem
	// 创建数组
	MakeArray
	// 创建hash
	MakeHash
	// 创建匿名函数
	MakeClosure
	// pop ?
	Pop
	// Noop ... ??
	Noop
	Add
	Sub
	Mul
	Div
	Mod
	Or
	And
	Not
	BitwiseOR
	BitwiseXOR
	BitwiseAND
	BitwiseNOT
	// ?
	LeftShift
	RightShift
	Equal
	NotEqual
	GreaterThan
	// GreaterThanEqual ...
	GreaterThanEqual
	Minus
	// jump if
	JumpIfFalse
	Jump
	// 调用函数
	Call
	// 返回
	Return
	ReturnValue
)

var definitions = map[Opcode]*Definition{
	LoadConstant:     {"LoadConstant", []int{2}},
	LoadBuiltin:      {"LoadBuiltin", []int{1}},
	AssignGlobal:     {"AssignGlobal", []int{2}},
	AssignLocal:      {"AssignLocal", []int{1}},
	LoadGlobal:       {"LoadGlobal", []int{2}},
	BindGlobal:       {"BindGlobal", []int{2}},
	LoadLocal:        {"LoadLocal", []int{1}},
	BindLocal:        {"BindLocal", []int{1}},
	LoadFree:         {"LoadFree", []int{1}},
	LoadModule:       {"LoadModule", []int{}},
	SetSelf:          {"SetSelf", []int{1}},
	LoadTrue:         {"LoadTrue", []int{}},
	LoadFalse:        {"LoadFalse", []int{}},
	LoadNull:         {"LoadNull", []int{}},
	GetItem:          {"GetItem", []int{}},
	SetItem:          {"SetItem", []int{}},
	MakeArray:        {"MakeArray", []int{2}},
	MakeHash:         {"MakeHash", []int{2}},
	MakeClosure:      {"MakeClosure", []int{2, 1}},
	Pop:              {"Pop", []int{}},
	Noop:             {"Noop", []int{}},
	Add:              {"Add", []int{}},
	Sub:              {"Sub", []int{}},
	Mul:              {"Mul", []int{}},
	Div:              {"Div", []int{}},
	Mod:              {"Mod", []int{}},
	Or:               {"Or", []int{}},
	And:              {"And", []int{}},
	Not:              {"Not", []int{}},
	BitwiseOR:        {"BitwiseOR", []int{}},
	BitwiseXOR:       {"BitwiseXOR", []int{}},
	BitwiseAND:       {"BitwiseAND", []int{}},
	BitwiseNOT:       {"BitwiseNOT", []int{}},
	LeftShift:        {"LeftShift", []int{}},
	RightShift:       {"RightShift", []int{}},
	Equal:            {"Equal", []int{}},
	NotEqual:         {"NotEqual", []int{}},
	GreaterThan:      {"GreaterThan", []int{}},
	GreaterThanEqual: {"GreaterThanEqual", []int{}},
	Minus:            {"Minus", []int{}},
	JumpIfFalse:      {"JumpIfFalse", []int{2}},
	Jump:             {"Jump", []int{2}},
	Call:             {"Call", []int{1}},
	Return:           {"Return", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

//[43]
//[4 0 0]
//[18 0 0 1]
func Make(op Opcode, operands ...int) []byte {
	// 操作码 => 对象
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}
	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
