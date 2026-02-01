#!/usr/bin/env bash
set -euo pipefail

GETH=./build/bin/geth

DATADIR1=/tmp/geth1
DATADIR2=/tmp/geth2
DATADIR3=/tmp/geth3
CLIQUE_DIR=/tmp/clique
GENESIS=$CLIQUE_DIR/genesis.json
PASSFILE=$CLIQUE_DIR/empty.pass

# ports
P2P1=30303; P2P2=30304; P2P3=30305
HTTP1=8545; HTTP2=8546; HTTP3=8547
AUTH1=8551; AUTH2=8552; AUTH3=8553

NETWORKID=1337
CHAINID=1337

LOG1=$CLIQUE_DIR/node1.log
LOG2=$CLIQUE_DIR/node2.log
LOG3=$CLIQUE_DIR/node3.log

echo "==> cleanup"
rm -rf "$DATADIR1" "$DATADIR2" "$DATADIR3" "$CLIQUE_DIR"
mkdir -p "$CLIQUE_DIR"
: > "$PASSFILE"

echo "==> create ONLY node1 account (empty password)"
OUT="$($GETH account new --datadir "$DATADIR1" --password "$PASSFILE")"
echo "$OUT"

ADDR="$(echo "$OUT" | awk '/Public address of the key/ {print $NF}' | tr -d '\r')"
if [[ -z "$ADDR" ]]; then
  ADDR="$(echo "$OUT" | awk '/Address:/ {gsub(/[{}]/,""); print $2}' | tr -d '\r')"
fi
if [[ -z "$ADDR" ]]; then
  echo "Failed to parse address"
  exit 1
fi
ADDR_LC="$(echo "$ADDR" | tr '[:upper:]' '[:lower:]')"
ADDR_NO0X="${ADDR_LC#0x}"

# clique extraData: 32B vanity + 20B signer + 65B sig
VANITY="$(python3 - <<'PY'
print("00"*32)
PY
)"
SIG="$(python3 - <<'PY'
print("00"*65)
PY
)"
EXTRADATA="0x${VANITY}${ADDR_NO0X}${SIG}"

echo "==> write genesis (signer = node1 only)"
cat > "$GENESIS" <<JSON
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
JSON

echo "==> init node1/node2/node3 with SAME genesis"
$GETH init --datadir "$DATADIR1" "$GENESIS" >/dev/null
$GETH init --datadir "$DATADIR2" "$GENESIS" >/dev/null
$GETH init --datadir "$DATADIR3" "$GENESIS" >/dev/null

echo "==> start 3 nodes in background (logs in $CLIQUE_DIR)"
# node1: unlock signer + feeRecipient
$GETH --datadir "$DATADIR1" --networkid "$NETWORKID" \
  --port "$P2P1" \
  --nat extip:127.0.0.1 \
  --http --http.addr 127.0.0.1 --http.port "$HTTP1" --http.api admin,eth,net,web3,txpool \
  --authrpc.addr 127.0.0.1 --authrpc.port "$AUTH1" \
  --unlock "$ADDR_LC" --password "$PASSFILE" \
  --miner.pending.feeRecipient "$ADDR_LC" \
  --nodiscover \
  --verbosity 5 >"$LOG1" 2>&1 &

PID1=$!

# node2/node3: plain peers
$GETH --datadir "$DATADIR2" --networkid "$NETWORKID" \
  --port "$P2P2" \
  --nat extip:127.0.0.1 \
  --http --http.addr 127.0.0.1 --http.port "$HTTP2" --http.api admin,eth,net,web3,txpool \
  --authrpc.addr 127.0.0.1 --authrpc.port "$AUTH2" \
  --nodiscover \
  --verbosity 5 >"$LOG2" 2>&1 &
PID2=$!

$GETH --datadir "$DATADIR3" --networkid "$NETWORKID" \
  --port "$P2P3" \
  --nat extip:127.0.0.1 \
  --http --http.addr 127.0.0.1 --http.port "$HTTP3" --http.api admin,eth,net,web3,txpool \
  --authrpc.addr 127.0.0.1 --authrpc.port "$AUTH3" \
  --nodiscover \
  --verbosity 5 >"$LOG3" 2>&1 &
PID3=$!

echo "PIDs: node1=$PID1 node2=$PID2 node3=$PID3"
echo "Signer address (node1): $ADDR_LC"
echo
echo "==> next: connect peers manually (since we’re local and deterministic)"
echo "1) get enode of node2:"
echo "   $GETH attach $DATADIR2/geth.ipc --exec admin.nodeInfo.enode"
echo "2) addPeer on node1 (paste enode):"
echo "   $GETH attach $DATADIR1/geth.ipc"
echo "   > admin.addPeer(\"enode://...\")"
echo "   > admin.peers.length"
echo "3) send tx batch on node1:"
echo "   > eth.accounts"
echo "   > eth.getBalance(eth.accounts[0])"
echo "   > for (let i=0;i<30;i++){eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1})}"
echo
echo "Logs:"
echo "  tail -f $LOG1"
echo "  tail -f $LOG2"
echo "  tail -f $LOG3"
