package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/mathetake/gasm/hostfunc"
	"github.com/mathetake/gasm/wasm"
	"github.com/stretchr/testify/require"
)

func testsIntAdd(a, b int64) int64 {
	return a + 2*b
}

func TestCallImports(t *testing.T) {
	builder := hostfunc.NewModuleBuilder()
	builder.MustSetFunction("guo", "intadd", func(machine *wasm.VirtualMachine) reflect.Value {
		return reflect.ValueOf(testsIntAdd)
	})
	ms := builder.Done()

	buf, err := ioutil.ReadFile("intadd.wasm")
	require.NoError(t, err)

	mod, err := wasm.DecodeModule(bytes.NewBuffer(buf))
	require.NoError(t, err)

	vm, err := wasm.NewVM(mod, ms)
	require.NoError(t, err)

	var a, b uint64 = 5, 3
	rets, retTypes, err := vm.ExecExportedFunction("plus10", a, b)
	require.NoError(t, err)
	t.Logf("returns: %v, returnTypes: %v", rets, retTypes)
	c := uint64(testsIntAdd(int64(a), int64(b))) + 10
	require.Equal(t, c, rets[0])
}

func TestCallStringConcat(t *testing.T) {
	builder := hostfunc.NewModuleBuilder()
	builder.MustSetFunction("guo", "concat", func(machine *wasm.VirtualMachine) reflect.Value {
		return reflect.ValueOf(testsConcat)
	})
	ms := builder.Done()

	buf, err := ioutil.ReadFile("strconcat.wasm")
	require.NoError(t, err)

	mod, err := wasm.DecodeModule(bytes.NewBuffer(buf))
	require.NoError(t, err)

	vm, err := wasm.NewVM(mod, ms)
	require.NoError(t, err)
	currentVM = vm

	rets, retTypes, err := vm.ExecExportedFunction("str_add")
	require.NoError(t, err)
	t.Logf("returns: %v, returnTypes: %v", rets, retTypes)
	result, err := vm.Mem.GetString(rets[0])
	if err != nil {
		t.Fatalf("GetString(%d) error: %v", rets[0], err)
	}
	t.Logf("returns: %s", result)
}

var currentVM *wasm.VirtualMachine

func testsConcat(p1, p2 uint64) uint64 {
	s1, err := currentVM.Mem.GetString(p1)
	if err != nil {
		panic(fmt.Sprintf("GetString(p1:%d) error: %v", p1, err))
	}
	s2, err := currentVM.Mem.GetString(p2)
	if err != nil {
		panic(fmt.Sprintf("GetString(p2:%d) error: %v", p2, err))
	}
	s := string(s1) + string(s2)
	r, err := currentVM.Mem.SetString(s)
	if err != nil {
		panic(fmt.Sprintf("SetString(%s) error: %v", s, err))
	}
	return r
}
