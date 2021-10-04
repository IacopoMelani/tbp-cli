package wallet

import "github.com/tyler-smith/go-bip39"

func GenerateNewMnemonic() (string, error) {

	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
