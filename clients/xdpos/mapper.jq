# Mapper from Geth genesis format to XDC-compatible genesis format.
# Uses ethash consensus so the node does not require XDPoS masternodes.

{
  config: {
    chainId: (.config.chainId // 1337),
    homesteadBlock: (.config.homesteadBlock // 0),
    eip150Block: (.config.eip150Block // 0),
    eip150Hash: (.config.eip150Hash // "0x0000000000000000000000000000000000000000000000000000000000000000"),
    eip155Block: (.config.eip155Block // 0),
    eip158Block: (.config.eip158Block // 0),
    byzantiumBlock: (.config.byzantiumBlock // 0),
    constantinopleBlock: (.config.constantinopleBlock // 0),
    petersburgBlock: (.config.petersburgBlock // 0),
    istanbulBlock: (.config.istanbulBlock // 0),
    berlinBlock: (.config.berlinBlock // 0),
    londonBlock: (.config.londonBlock // 0),
    ethash: {}
  },
  nonce: "0x0",
  timestamp: (.timestamp // "0x0"),
  extraData: "0x0000000000000000000000000000000000000000000000000000000000000000",
  gasLimit: (.gasLimit // "0x47b760"),
  difficulty: (.difficulty // "0x1"),
  mixHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
  coinbase: "0x0000000000000000000000000000000000000000",
  alloc: (.alloc // {}),
  number: (.number // "0x0"),
  gasUsed: (.gasUsed // "0x0"),
  parentHash: (.parentHash // "0x0000000000000000000000000000000000000000000000000000000000000000")
}
