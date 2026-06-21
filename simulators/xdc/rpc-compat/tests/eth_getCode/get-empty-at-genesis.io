// eth_getCode for a non-contract account at genesis returns empty.
>> {"jsonrpc":"2.0","id":1,"method":"eth_getCode","params":["0x065551f0dcac6f00cae11192d462db709be3758c","0x0"]}
<< {"jsonrpc":"2.0","id":1,"result":"0x"}
