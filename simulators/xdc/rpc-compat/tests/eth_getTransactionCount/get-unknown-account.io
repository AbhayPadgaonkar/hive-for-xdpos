// eth_getTransactionCount for an unknown account returns zero.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getTransactionCount","params":["0x000000000000000000000000000000000000dead","latest"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x0"}
