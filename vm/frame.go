package vm

import (
	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/object"
)

type Frame struct {
	cl          *object.Closure			// 回调
	ip          int						// 指令指针，Instruction Pointer
	basePointer int						// 基础指针
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}

	return f
}

// NextOp ...
func (f *Frame) NextOp() code.Opcode {
	return code.Opcode(f.Instructions()[f.ip+1])
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
