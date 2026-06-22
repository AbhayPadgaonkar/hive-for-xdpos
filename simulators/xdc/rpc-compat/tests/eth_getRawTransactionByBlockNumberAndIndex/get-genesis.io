// eth_getRawTransactionByBlockNumberAndIndex returns empty for a genesis block without transactions.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getRawTransactionByBlockNumberAndIndex","params":["0x0","0x0"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x"}
