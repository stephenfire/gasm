package tests

import (
	"bytes"
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

func testsGuoStrslength(apointer, bpointer int32) int32 {

}
