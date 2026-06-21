// eth_getTransactionByBlockNumberAndIndex returns null for genesis block with no transactions.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByBlockNumberAndIndex","params":["0x0","0x0"]}
<< {"jsonrpc":"2.0","id":1,"result":null}
