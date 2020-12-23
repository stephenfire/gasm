package wasm

import (
	"fmt"
	"math"
	"reflect"
)

type (
	VirtualMachineFunction interface {
		Call(vm *VirtualMachine)
		FunctionType() *FunctionType
	}
	HostFunction struct {
		ClosureGenerator func(vm *VirtualMachine) reflect.Value
		function         reflect.Value // should be set at the time of VM creation
		Signature        *FunctionType
	}
	NativeFunction struct {
		Signature *FunctionType
		NumLocal  uint32
		Body      []byte
		Blocks    map[uint64]*NativeFunctionBlock
	}
	NativeFunctionBlock struct {
		StartAt, ElseAt, EndAt uint64
		BlockType              *FunctionType
		BlockTypeBytes         uint64
	}
)

var (
	_ VirtualMachineFunction = &HostFunction{}
	_ VirtualMachineFunction = &NativeFunction{}
)

func (h *HostFunction) FunctionType() *FunctionType {
	return h.Signature
}

func (n *NativeFunction) FunctionType() *FunctionType {
	return n.Signature
}

func (h *HostFunction) Call(vm *VirtualMachine) {
	tp := h.function.Type()
	in := make([]reflect.Value, tp.NumIn())
	for i := len(in) - 1; i >= 0; i-- {
		val := reflect.New(tp.In(i)).Elem()
		raw := vm.OperandStack.Pop()
		kind := tp.In(i).Kind()

		switch kind {
		case reflect.Float64, reflect.Float32:
			val.SetFloat(math.Float64frombits(raw))
		case reflect.Uint32, reflect.Uint64:
			val.SetUint(raw)
		case reflect.Int32, reflect.Int64:
			val.SetInt(int64(raw))
		case reflect.String:
			bs, err := vm.Mem.GetString(raw)
			if err != nil {
				panic(fmt.Sprintf("GetString(%d) error: %v", raw, err))
			}
			val.SetString(string(bs))
			// ptr := int(raw)
			// if ptr >= len(vm.Mem.Memory) {
			// 	panic("memory out of bound")
			// }
			// strLen := 0
			// p := ptr
			// for ; p < len(vm.Mem.Memory) && vm.Mem.Memory[p] != 0x0; p++ {
			// 	strLen++
			// }
			// if p == len(vm.Mem.Memory) {
			// 	panic("invalid memory index")
			// }
			// bs := make([]byte, strLen)
			// if strLen > 0 {
			// 	copy(bs, vm.Mem.Memory[ptr:p])
			// }
			// val.SetString(string(bs))
		default:
			panic("invalid input type")
		}
		in[i] = val
	}

	for _, ret := range h.function.Call(in) {
		switch ret.Kind() {
		case reflect.Float64, reflect.Float32:
			vm.OperandStack.Push(math.Float64bits(ret.Float()))
		case reflect.Uint32, reflect.Uint64:
			vm.OperandStack.Push(ret.Uint())
		case reflect.Int32, reflect.Int64:
			vm.OperandStack.Push(uint64(ret.Int()))
		case reflect.String:

		default:
			panic("invalid return type")
		}
	}
}

func (n *NativeFunction) Call(vm *VirtualMachine) {
	al := len(n.Signature.InputTypes)
	locals := make([]uint64, n.NumLocal+uint32(al))
	for i := 0; i < al; i++ {
		locals[al-1-i] = vm.OperandStack.Pop()
	}

	prev := vm.ActiveContext
	vm.ActiveContext = &NativeFunctionContext{
		Function:   n,
		Locals:     locals,
		LabelStack: NewVirtualMachineLabelStack(),
	}
	vm.execNativeFunction()
	vm.ActiveContext = prev
}

func (vm *VirtualMachine) execNativeFunction() {
	for ; int(vm.ActiveContext.PC) < len(vm.ActiveContext.Function.Body); vm.ActiveContext.PC++ {
		switch op := vm.ActiveContext.Function.Body[vm.ActiveContext.PC]; OptCode(op) {
		case OptCodeReturn:
			return
		default:
			virtualMachineInstructions[op](vm)
		}
	}
}
