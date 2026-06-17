// eth_getCode for an unknown account returns empty.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getCode","params":["0x0000000000000000000000000000000000000000","latest"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x"}
