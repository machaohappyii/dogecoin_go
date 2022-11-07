package main

import "go_btc/dogeCoin"

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
