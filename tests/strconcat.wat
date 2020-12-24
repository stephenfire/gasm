(module
    (import "guo" "concat" (func $host_concat (param i64) (param i64) (result i64)))
    (memory 1)
    (data (offset (i32.const 10)) "this is")
    (data (offset (i32.const 100)) " big mouth")
    (data (offset (i32.const 200)) "!!!")
    (func (export "str_add") (result i64)
        i64.const 10
        i64.const 100
        call $host_concat
        i64.const 200
        call $host_concat
    )
)