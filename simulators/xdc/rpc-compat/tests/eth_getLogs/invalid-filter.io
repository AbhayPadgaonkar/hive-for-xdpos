// eth_getLogs with invalid filter returns error.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getLogs","params":["invalid"]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":""}}
