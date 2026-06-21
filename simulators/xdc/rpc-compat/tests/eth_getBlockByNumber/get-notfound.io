// eth_getBlockByNumber returns null for a block that does not exist.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["0xffff",false]}
<< {"jsonrpc":"2.0","id":1,"result":null}
