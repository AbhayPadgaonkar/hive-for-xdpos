// eth_signTransaction fails because there is no unlocked account over HTTP.
>> {"jsonrpc":"2.0","id":1,"method":"eth_signTransaction","params":[{"from":"0x746249c61f5832c5eed53172776b460491bdcd5c","to":"0x000000000000000000000000000000000000dead","value":"0x1","gas":"0x5208","gasPrice":"0x430e23400","nonce":"0x0"}]}
<< {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":""}}
