package dogeCoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type balanceResource struct {
	Balance string `json:"balance"`
	Success int    `json:"success"`
}

type utxoResource struct {
	UnspentOutputs []struct {
		TxHash        string `json:"tx_hash"`
		TxOutputN     int    `json:"tx_output_n"`
		Script        string `json:"script"`
		Value         string `json:"value"`
		Confirmations int    `json:"confirmations"`
		Address       string `json:"address"`
	} `json:"unspent_outputs"`
	Success int `json:"success"`
}

type txData struct {
	Tx string `json:"tx"`
}

type sendResource struct {
	TxHash  string `json:"tx_hash"`
	Success int    `json:"success"`
}

// api 获取余额
func GetBalance(address string) (relBalance balanceResource, err error) {
	url := "https://dogechain.info/api/v1/address/balance/" + address
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("获取数据err", err)
		return relBalance, err
	}
	if err := json.Unmarshal(body, &relBalance); err != nil {
		fmt.Println("json解析错误:", err)
		return relBalance, err
	}
	fmt.Printf("%+v", relBalance)
	return relBalance, err
}

// api获取可以转账的utxo
func GetUtxo(address string) (relUtxo utxoResource, err error) {
	url := "https://dogechain.info/api/v1/unspent/" + address
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("获取数据err", err)
		return relUtxo, err
	}
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &relUtxo); err != nil {
		fmt.Println("json解析错误:", err)
		return relUtxo, err
	}
	return relUtxo, err
}

// apt 广播交易
func SendTx(tx string) (relSend sendResource, err error) {
	data := txData{Tx: tx}
	txData, _ := json.Marshal(&data)
	reader := bytes.NewReader(txData)
	url := "https://dogechain.info/api/v1/pushtx"
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("获取数据err", err)
		return relSend, err
	}
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &relSend); err != nil {
		fmt.Println("json解析错误:", err)
		return relSend, err
	}
	fmt.Printf("sendTx result=:%+v", relSend)
	return relSend, err
}
