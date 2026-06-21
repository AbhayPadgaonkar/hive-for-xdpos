// eth_call to a non-existent account returns empty result.
>> {"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"to":"0x000000000000000000000000000000000000dead","data":"0x"},"latest"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x"}
