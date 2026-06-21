// eth_getTransactionByHash returns null for a missing transaction hash.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByHash","params":["0x0000000000000000000000000000000000000000000000000000000000000001"]}
<< {"jsonrpc":"2.0","id":1,"result":null}
