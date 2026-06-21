// eth_getBalance with invalid address format returns error.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0xnotanaddress","latest"]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":""}}
