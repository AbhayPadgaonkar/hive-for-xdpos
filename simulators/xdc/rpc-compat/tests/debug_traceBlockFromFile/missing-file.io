// debug_traceBlockFromFile returns an error for a missing file.
>> {"jsonrpc":"2.0","id":1,"method":"debug_traceBlockFromFile","params":["/tmp/nonexistent.rlp",{}]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":""}}
