package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/mr-tron/base58"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func PkToAddress(p string) string {
	privateBytes, err := hex.DecodeString(p)
	if err != nil {
		fmt.Println("解析私钥错误")
	}
	privateKey := crypto.ToECDSAUnsafe(privateBytes)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//if 0 > 1 {
	//	fmt.Println("publicKey:", hexutil.Encode(publicKeyBytes)[2:])
	//}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address = "41" + address[2:]
	//fmt.Println("address hex: ", address)
	addb, _ := hex.DecodeString(address)
	hash1 := s256(s256(addb))
	secret := hash1[:4]
	for _, v := range secret {
		addb = append(addb, v)
	}
	return base58.Encode(addb)
}

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func GetBaseAddress(address string) string {
	if len(address) > 0 {
		address = "41" + address[2:]
		addb, _ := hex.DecodeString(address)
		hash1 := s256(s256(addb))
		secret := hash1[:4]
		for _, v := range secret {
			addb = append(addb, v)
		}
		return base58.Encode(addb)
	}
	return ""
}

func RemoveIndex(slice []string, index int) []string {
	tmp := make([]string, 0, len(slice))
	for num, v := range slice {
		if num != index {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

func Debug(error string) {
	//logS.LogError().Println(error)
	//fmt.Println(error)
}

func HexToBigint(s string) float64 {
	n := new(big.Int)
	n, _ = n.SetString(s[2:], 16)
	number, error := strconv.ParseFloat(n.String(), 64)
	if error != nil {
		fmt.Println(number)
	}
	return number
}

func RealCheckAllowance(cli *client.GrpcClient, token, addrA, addrB string) (error, bool, int64) {
	addrOne, err := address.Base58ToAddress(addrA)
	addrTwo, err := address.Base58ToAddress(addrB)
	if err != nil {
		return errors.New("network error"), false, 0
	}
	req := "0xdd62ed3e" + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrOne.Hex())-2:] + addrOne.Hex()[2:]
	req += "0000000000000000000000000000000000000000000000000000000000000000"[len(addrTwo.Hex())-2:] + addrTwo.Hex()[2:]
	result, err := cli.TRC20Call("", token, req, true, 0)
	if err != nil {
		fmt.Println(err)
		return errors.New("network error"), false, 0
	}
	data := common.BytesToHexString(result.GetConstantResult()[0])
	r, err := cli.ParseTRC20NumericProperty(data)
	if err != nil {
		fmt.Println(err)
		return errors.New("network error"), false, 0
	}
	if r == nil {
		return errors.New("network error"), false, 0
	}
	fmt.Println("数据")
	return nil, !(r.Int64() == 0), r.Int64()
}

func TransferFromCoin(cli *client.GrpcClient, pk, token, addrA, addrB string, amount *big.Int) (error, string) {
	addrOne, err := address.Base58ToAddress(addrA)
	addrTwo, err := address.Base58ToAddress(addrB)
	if err != nil {
		return err, ""
	}
	req := "0x23b872dd" + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrOne.Hex())-2:] + addrOne.Hex()[2:]
	req += "0000000000000000000000000000000000000000000000000000000000000000"[len(addrTwo.Hex())-2:] + addrTwo.Hex()[2:]
	ab := common.LeftPadBytes(amount.Bytes(), 32)
	req += common.Bytes2Hex(ab)
	tx, err := cli.TRC20Call(PkToAddress(pk), token, req, true, 500000)
	signTx, err := sign.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return err, ""
	}
	_, err = cli.Broadcast(signTx)
	if err != nil {
		return err, ""
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	return nil, txid
}
