# 更改源代码的地方 且非 `p2p/costspace` 文件夹下的

使用了下面的字符串进行标识。

```go
//-------test hyperbolic cost model-------
//-------end test hyperbolic cost model-------
```

第一行为开始，第二行为结束。

## 本地测试

### 创造网络

`./p2p/costspace/start-hyperbolic-network.sh`

### node2js console

```shell
./build/bin/geth attach /tmp/geth2/geth.ipc
```

```js
admin.nodeInfo.enode
```

### node3js console

```shell
./build/bin/geth attach /tmp/geth3/geth.ipc
```

```js
admin.nodeInfo.enode
```

### node1 js console 

```shell
./build/bin/geth attach /tmp/geth1/geth.ipc
```

```shell
eth.accounts
```

```shell
eth.getBalance(eth.accounts[0])
```

```js
admin.addPeer("")
admin.peers.length
```

make tx

```js
for (let i = 0; i < 30; i++) {eth.sendTransaction({from: eth.accounts[0],to: eth.accounts[0],value: 1})}
```

stop the network

`./p2p/costspace/stop-network.sh`

## formate of the genesis.json

```json
{
  "config": {
    "chainId": $CHAINID,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,

    "terminalTotalDifficulty": 0,
    "shanghaiTime": 0,
    "cancunTime": 0,
    "pragueTime": 0,
    "osakaTime": 0,

    "clique": { "period": 2, "epoch": 30000 },
    "blobSchedule": {
      "cancun": {
        "target": 3,
        "max": 6,
        "baseFeeUpdateFraction": 3338477
      },
      "prague": {
        "target": 3,
        "max": 6,
        "baseFeeUpdateFraction": 3338477
      },
      "osaka": {
        "target": 3,
        "max": 6,
        "baseFeeUpdateFraction": 3338477
      }
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x2fefd8",
  "extraData": "$EXTRADATA",
  "alloc": {
    "$ADDR_LC": { "balance": "0x3635C9ADC5DEA00000" }
  }
}
```

