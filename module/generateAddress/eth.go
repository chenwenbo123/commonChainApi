package generate

import (
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func ProduceEth() *SingleData {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)

	seed := bip39.NewSeed(mnemonic, "") //这里可以选择传入指定密码或者空字符串，不同密码生成的助记词不同

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/195'/0'/0/0") //最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, _ := wallet.PrivateKeyHex(account)

	return &SingleData{
		Mnemonic:   mnemonic,
		PrivateKey: "0x" + privateKey,
		Address:    account.Address.Hex(),
	}
}
