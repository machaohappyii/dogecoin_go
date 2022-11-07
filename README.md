# go dogecoin demo

[![Documentation (master)](https://img.shields.io/badge/docs-master-59f)](https://github.com/coming-chat/wallet-SDK)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://github.com/aptos-labs/aptos-core/blob/main/LICENSE)

golang 版本的dogecoin链上的操作

## Clone

```sh
git clone https://github.com/machaohappyii/dogecoin_go.git
cd go_dogecoin
go mod tidy
```

## Usage
```sh
doge浏览器:
https://dogechain.info/

doge官方api:
https://dogechain.info/api/blockchain_api

go-hdwallet库生成doge账户:
https://github.com/foxnut/go-hdwallet

btcsuite库进行构建交易:  
https://github.com/btcsuite  
```

## Function

- [x] 创建账号
- [x] 查询余额
- [x] UTXO
- [x] 转账

### Account

```go
func main() {
    //生成doge的账户
    dogeCoin.CreateAccount()
    
    //获取账户余额
    dogeCoin.Balance()
    
    //获取用户的utxo
    dogeCoin.Utxo()
    
    //转账
    dogeCoin.Transfer()
}
```


