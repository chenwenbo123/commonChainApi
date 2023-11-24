package tron

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	goClient "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/gogo/protobuf/proto"
)

func SignTransaction(transaction *core.Transaction, privateKey string) (*core.Transaction, error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode private key error: %v", err)
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	defer zeroKey(priv)
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}
	fmt.Println(rawData)
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, priv)
	if err != nil {
		return nil, fmt.Errorf("sign error: %v", err)
	}
	transaction.Signature = append(transaction.Signature, signature)
	return transaction, nil
}

// zeroKey zeroes a private key in memory.
func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

// UpdateAccountPermission change account permission
// func (g *GrpcClient) UpdateAccountPermission(from string, owner, witness map[string]interface{}, actives []map[string]interface{}) (*api.TransactionExtention, error) {
func UpdateAccountPermission(g *goClient.GrpcClient, from string, owner *core.Permission, witness map[string]interface{}, actives []*core.Permission) (*api.TransactionExtention, error) {
	if len(actives) > 8 {
		return nil, fmt.Errorf("cant have more than 8 active operations")
	}
	contract := &core.AccountPermissionUpdateContract{
		Owner: owner,
	}
	var Z1 error
	contract.OwnerAddress, Z1 = common.DecodeCheck(from)
	if Z1 != nil {

	}
	contract.Actives = actives
	if witness != nil {
	}
	//ctx, cancel := g.getContext()
	ctx := context.Background()
	tx, err := g.Client.AccountPermissionUpdate(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}
