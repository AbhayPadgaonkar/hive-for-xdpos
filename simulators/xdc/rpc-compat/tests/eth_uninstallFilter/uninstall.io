// eth_uninstallFilter returns false for a non-existent filter id.
>> {"jsonrpc":"2.0","id":1,"method":"eth_uninstallFilter","params":["0x0000000000000000000000000000000000000000000000000000000000000000"]}
<< {"jsonrpc":"2.0","id":1,"result":false}
