// eth_getStorageAt coerces an invalid storage key to zero on this XDC build.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getStorageAt","params":["0x000000000000000000000000000000000000dead","latest","latest"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}
