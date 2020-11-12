package examples

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mathetake/gasm/wasi"
	"github.com/mathetake/gasm/wasm"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	buf, err := ioutil.ReadFile("wasm/hello.wasm")
	require.NoError(t, err)

	mod, err := wasm.DecodeModule(bytes.NewBuffer(buf))
	require.NoError(t, err)

	vm, err := wasm.NewVM(mod, wasi.Modules)
	require.NoError(t, err)

	ret, retTypes, err := vm.ExecExportedFunction("main")
	fmt.Printf("%d %s %v", ret, retTypes, err)
	//
	// for _, c := range []struct {
	// 	in, exp int32
	// }{
	// 	{in: 20, exp: 6765},
	// 	{in: 10, exp: 55},
	// 	{in: 5, exp: 5},
	// } {
	// 	ret, retTypes, err := vm.ExecExportedFunction("fibonacci", uint64(c.in))
	// 	require.NoError(t, err)
	// 	require.Len(t, ret, len(retTypes))
	// 	require.Equal(t, wasm.ValueTypeI32, retTypes[0])
	// 	require.Equal(t, c.exp, int32(ret[0]))
	// }
	//
}
