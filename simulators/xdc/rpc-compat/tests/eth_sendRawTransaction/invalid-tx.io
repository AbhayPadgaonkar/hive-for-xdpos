// eth_sendRawTransaction rejects an empty transaction.
>> {"jsonrpc":"2.0","id":1,"method":"eth_sendRawTransaction","params":["0x"]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":""}}
