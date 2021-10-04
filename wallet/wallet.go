package wallet

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func NewWalletFromMnemonic(mnemonic string) (*hdwallet.Wallet, error) {

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func GetAccountFromWallet(wallet *hdwallet.Wallet) (accounts.Account, error) {

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		return accounts.Account{}, err
	}

	return account, nil
}

func GenerateKeyFromMnemonicAndStore(filename string, auth string, mnemonic string) error {

	wallet, err := NewWalletFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	account, err := GetAccountFromWallet(wallet)
	if err != nil {
		return err
	}

	privKey, err := wallet.PrivateKey(account)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	key := &keystore.Key{
		Address:    account.Address,
		PrivateKey: privKey,
		Id:         id,
	}

	jsonKey, err := keystore.EncryptKey(key, auth, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, jsonKey, 0600); err != nil {
		return err
	}

	return nil
}


	// Generate a mnemonic for memorization or user-friendly seeds
	// entropy, _ := bip39.NewEntropy(256)
	// mnemonic, _ := bip39.NewMnemonic(entropy)

	// // Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	// seed := bip39.NewSeed(mnemonic, "")

	// masterKey, _ := bip32.NewMasterKey(seed)
	// publicKey := masterKey.PublicKey()

	// // Display mnemonic and keys
	// fmt.Println("Mnemonic: ", mnemonic)
	// fmt.Println("Master private key: ", masterKey)
	// fmt.Println("Master public key: ", publicKey)

	// wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	// account, err := wallet.Derive(path, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// privKey, err := wallet.PrivateKey(account)
	// if err != nil {
	// 	panic(err)
	// }

	// id, err := uuid.NewRandom()
	// if err != nil {
	// 	panic(err)
	// }

	// key := &keystore.Key{
	// 	Address:    account.Address,
	// 	PrivateKey: privKey,
	// 	Id:         id,
	// }

	// jsonKey, err := keystore.EncryptKey(key, "123", keystore.StandardScryptN, keystore.StandardScryptP)
	// if err != nil {
	// 	panic(err)
	// }

	// os.WriteFile("key", jsonKey, 0600)

	// fmt.Println(account.Address.Hex())