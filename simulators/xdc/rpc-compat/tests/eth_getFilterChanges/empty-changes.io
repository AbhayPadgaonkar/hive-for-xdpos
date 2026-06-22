// eth_getFilterChanges returns an error for a non-existent filter id.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getFilterChanges","params":["0x0000000000000000000000000000000000000000000000000000000000000000"]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":""}}
