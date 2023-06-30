# FVM Deal Making demo

This repo is a demo on how to make deals using a FVM smart contract.

It's based on [fvm-starter-kit-deal-making](https://github.com/filecoin-project/fvm-starter-kit-deal-making).

We have a Solidity contract [DealClient.sol](./contracts/DealClient.sol) responsible for making deals and a Go CLI for interacting with that contract.

## Deploying contract

We already have a contract deployed in Calibration network at address [0x5b618b1Ec92d7Df96307c88b691Ad42913e60e96](https://calibration.filfox.info/en/address/0x5b618b1Ec92d7Df96307c88b691Ad42913e60e96) that you can use to make deals. If you want to deploy your own contract, the command is:

```bash
PRIVATE_KEY='FILL WITH PRIVATE KEY' npx hardhat run --network [NETWORK] scripts/deploy.js
```

If you want to test the contract on a local network you can make use of [filecoin-fvm-localnet](https://github.com/filecoin-project/filecoin-fvm-localnet) to have local Filecoin network running in your computer.

## Creating a deal

Before creating a deal for a file you want to store in Filecoin, you have to prepare your data by generating a CAR file.

The easiest way is to use https://data.lighthouse.storage/.

Once you have your data prepared you can make a deal by using the following command:

```bash
cd go-dealmaker

go run *.go create --rpc-endpoint https://api.calibration.node.glif.io/ \
    --contract 0x8B9B28A202d91FeAC7dC63840A743b17B8F70718 \
    --piece-cid baga6ea4seaqohrguieio3w2e23p5iaktc24hjjlkw2dthp4g2q4vimci7jvbwgq  \
    --piece-size 4194304 \
    --verified  \
    --payload-cid bafybeidgfkqzkcjfszxozssw7mhtnbg4v2lj6o2lydugfqbeovgak4zftq \
    --start-epoch 690697 \
    --end-epoch 1690697 \
    --location-ref "https://data-depot.lighthouse.storage/api/download/download_car?fileId=397cc6e3-9c46-4974-9b82-4026a56548ac.car" \
    --car-size 4050510 \
    --private-key [PRIVATE KEY] \
    --chain-id 314159
```

[Here](https://github.com/filecoin-project/community/discussions/634) you can find some guidelines on how to set a proper `start-epoch`. Usually, current epoch + 30000 is a good heuristic.

It takes roughly 1h for the Store Provider to publish that event.

## Checking the deal status

```bash
cd go-dealmaker 

go run *.go status --rpc-endpoint https://api.calibration.node.glif.io/ \
    --contract 0x8B9B28A202d91FeAC7dC63840A743b17B8F70718 \
    --piece-cid baga6ea4seaqohrguieio3w2e23p5iaktc24hjjlkw2dthp4g2q4vimci7jvbwgq
```
