(module
    (import "guo" "intadd" (func $host_intadd (param i64) (param i64) (result i64)))
    (memory 1)
    (func (export "plus10") (param i64) (param i64) (result i64)
        get_local 0
        get_local 1
        call $host_intadd
        i64.const 10
        i64.add
    )
)