package main

import (
	"fmt"

	"go.dedis.ch/dela/crypto"
	"go.dedis.ch/dela/crypto/bls"
)

func main() {
	// Message To Sign
	message := []byte("Hello Bro")
	// User 1 Or Signer 1, Note Here We Have Created Random User Key Pair Or You Can Say Account, You Can Generate Public Key Using 256 Bit Private Key.
	signer1 := bls.NewSigner()
	// User 2 Or Signer 2
	signer2 := bls.NewSigner()

	// Signature Signing By User 1
	signer1Signature, err := signer1.Sign(message)
	if err != nil {
		fmt.Println("Something Went Wrong During Signature Signing By User1: ", err)
	}

	// Signature Signing By User 1
	signer2Signature, err := signer2.Sign(message)
	if err != nil {
		fmt.Println("Something Went Wrong During Signature Signing By User2: ", err)
	}

	// Combining Or Aggregating Signatures Signed By Both Users
	combinedSignature, err := signer1.Aggregate(signer1Signature, signer2Signature)
	if err != nil {
		fmt.Println("Something Went Wrong During Combining Signature Of Both Users: ", err)
	}

	// Generating Verifier To Verify Combined Signature.
	verifier, err := signer1.GetVerifierFactory().FromArray([]crypto.PublicKey{signer1.GetPublicKey(), signer2.GetPublicKey()})
	if err != nil {
		fmt.Println("Something Went Wrong During Generating Verifier: ", err)
	}

	// Verifying Aggregate Signature
	err = verifier.Verify(message, combinedSignature)
	if err != nil {
		fmt.Println("Signature Verification Failed: ", err)
	} else {
		fmt.Println("Signature Verification Done")
	}
}
