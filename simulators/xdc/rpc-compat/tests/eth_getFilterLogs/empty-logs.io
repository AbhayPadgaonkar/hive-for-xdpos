// eth_getFilterLogs returns an error for a non-existent filter id.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getFilterLogs","params":["0x0000000000000000000000000000000000000000000000000000000000000000"]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":""}}
