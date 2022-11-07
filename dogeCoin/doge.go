package dogeCoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/foxnut/go-hdwallet"
	"math/big"
	"strconv"
)

// 主要对象btc主网络 修改 网络配置信息 @hunter
var DOGEParams = chaincfg.MainNetParams

// 主要点  就是修改doge的网络 @hunter
func setDogeConfig() {
	DOGEParams.PubKeyHashAddrID = 0x1e // 30
	DOGEParams.ScriptHashAddrID = 0x16 // 22
	DOGEParams.PrivateKeyID = 0x9e     // 158
	DOGEParams.HDCoinType = 0
}

// 创建一个账户
func CreateAccount() {
	var mnemonic = "range sheriff try enroll deer over ten level bring display stamp recycle"
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		panic(err)
	}
	//路径  m/44'/3'/0'/0/0
	var index uint32 = 1
	wallet, _ := master.GetWallet(hdwallet.CoinType(hdwallet.DOGE), hdwallet.AddressIndex(index))
	address, _ := wallet.GetAddress()
	Private, err := wallet.GetKey().PrivateWIF(true)
	Publick := wallet.GetKey().PublicHex(true)
	if err != nil {
		panic(err)
	}
	fmt.Println("BTC地址：", address)
	fmt.Println("BTC私钥：", Private)
	fmt.Println("BTC公钥：", Publick)
}

// 获取自己的doge余额
func Balance() {
	address := "DJrV3roiSDXk2r7nCesVPd5asYcb7qAwxV"
	rel, _ := GetBalance(address)
	if rel.Success == 1 {
		fmt.Println(rel.Balance)
	} else {
		fmt.Println("balance err")
	}
}

// 获取自己可以转账的UTXO
func Utxo() {
	address := "DJrV3roiSDXk2r7nCesVPd5asYcb7qAwxV"
	rel, _ := GetUtxo(address)
	if rel.Success == 1 {
		fmt.Printf("%+v", rel.UnspentOutputs)
	} else {
		fmt.Println("Invalid address")
	}
}

//根据私钥获取地址
func getAddress(privKey string) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}
	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &DOGEParams)
	address := addrPubKey.EncodeAddress()
	fmt.Println("from", address)
	return address, nil
}

// doge 转账
func Transfer() {
	//修改doge的网络--主要是这个 @hunter 注意看下
	setDogeConfig()
	//转账人的私钥
	fromPrivKey := "QS4RDvPG9RgZTpNovgX4KDCjio4BWttE8WN6ir1oh1fGPzHuwECV"
	//转账人地址 也是找零地址
	from, err := getAddress(fromPrivKey)

	if err != nil {
		fmt.Println("getAddress fail", err)
		return
	}
	//转账人
	to := "DHLA3rJcCjG2tQwvnmoJzD5Ej7dBTQqhHK"
	//转多少
	var toVal int64 = 10000000
	//固定gas费 最稳妥的0.05    0.01也可以
	var gasFee int64 = 5000000
	// 获取构建交易
	tx, err := CreateTx(fromPrivKey, from, to, toVal, gasFee)
	if err != nil {
		fmt.Println("err trsnfer fail", err)
		return
	}
	fmt.Println(tx)
	// 发送交易
	rel, err := SendTx(tx)
	if err != nil {
		fmt.Println("err SendTx fail", err)
		return
	}
	if rel.Success == 1 {
		fmt.Println("SendTx success", rel.Success)
		hash := rel.TxHash
		fmt.Println("send success hash", hash)
	} else {
		fmt.Println("err SendTx fail", err)
		return
	}
}

// 创建交易
// 拿到自己utxo 加入 input 公钥验证  output构建转账人 找零地址
func CreateTx(privKey string, from string, to string, amount int64, fee int64) (string, error) {
	utxo, err := GetUtxo(from)
	if err != nil {
		return "", err
	}

	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// 总的utxo余额
	var utxoToal = big.NewInt(0)

	// 构建input
	for _, list := range utxo.UnspentOutputs {
		val, _ := strconv.ParseInt(list.Value, 10, 64)
		utxoToal.Add(utxoToal, big.NewInt(val))
		txid := list.TxHash
		txidIndex := list.TxOutputN
		utxoHash, err := chainhash.NewHashFromStr(txid)
		if err != nil {
			return "", err
		}
		fmt.Println("txid", txid)
		fmt.Println("txidIndex", txidIndex)
		outPoint := wire.NewOutPoint(utxoHash, uint32(txidIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		redeemTx.AddTxIn(txIn)
	}

	// 构建转账用户
	toAddr, err := btcutil.DecodeAddress(to, &DOGEParams)
	if err != nil {
		return "", err
	}
	toAddrByte, err := txscript.PayToAddrScript(toAddr)
	if err != nil {
		return "", err
	}
	redeemTxOut := wire.NewTxOut(int64(amount), toAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	totalAmount := big.NewInt(0).Add(big.NewInt(amount), big.NewInt(fee))
	fmt.Println(utxoToal, totalAmount)

	//找零余额 utxo剩下的钱 还给本人
	changeVal := utxoToal.Sub(utxoToal, totalAmount)
	fmt.Println("reedVal", changeVal, amount, fee, changeVal.Cmp(big.NewInt(0)))

	// 说明钱不够
	if changeVal.Cmp(big.NewInt(0)) == -1 {
		return "", fmt.Errorf("balance 不错")
	}

	if changeVal.Cmp(big.NewInt(0)) > 0 {
		// 构建找零钱用户
		fromAddr, err := btcutil.DecodeAddress(from, &DOGEParams)
		if err != nil {
			return "", err
		}
		fromAddrByte, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			return "", err
		}
		//转成int64
		changeVal, _ := strconv.ParseInt(changeVal.String(), 10, 64)
		fromTxOut := wire.NewTxOut(changeVal, fromAddrByte)
		redeemTx.AddTxOut(fromTxOut)
	}
	//根据私钥获取签名公钥
	pkScript, err := pubKeyScript(from)
	if err != nil {
		return "", err
	}
	// 对 input 签名
	finalRawTx, err := SignTx(privKey, pkScript, redeemTx, utxo)
	fmt.Println("finalRawTx", finalRawTx)
	return finalRawTx, nil

}

// 获取签名公钥
func pubKeyScript(from string) (string, error) {
	fromAddr, err := btcutil.DecodeAddress(from, &DOGEParams)
	if err != nil {
		panic(err)
	}
	fromAddrByte, err := txscript.PayToAddrScript(fromAddr)
	pubKeyScript := hex.EncodeToString(fromAddrByte)
	return pubKeyScript, err
}

// 给 每个 input 签名
func SignTx(privKey string, pkScript string, redeemTx *wire.MsgTx, utxo utxoResource) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return "", nil
	}

	for index, _ := range utxo.UnspentOutputs {
		signature, err := txscript.SignatureScript(redeemTx, index, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			return "", nil
		}
		redeemTx.TxIn[index].SignatureScript = signature
	}
	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
