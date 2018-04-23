package main

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
)

type Patient struct {
	HosID string
	PatientID string
}

func main() {

	var m = make(map[Patient]*rsa.PrivateKey) // a map storing

	var privatekey *rsa.PrivateKey
	var err error
	reader := rand.Reader
	bitSize := 2048

	a := []Patient{newPatient("pubhos1", "cabi"), newPatient("pubhos1", "aa"), newPatient("pubhos1", "siuon")}

	for i := 0; i < len(a); i++ {
		privatekey, err = rsa.GenerateKey(reader, bitSize)
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
		}
		if _, ok := m[a[i]]; ok {
			// do something
			fmt.Println("key existed", a[i])
		}	else {
			m[a[i]] = privatekey
		}
	}

	for i := 0; i < len(a); i++ {
		//fmt.Println(m[a[i]].PublicKey)
	}

	originalText := "encrypt this golang"
	fmt.Println(originalText)
	var mapToAES = make(map[Patient][]string) // (PID, HOSID) -> encEDCSA(PublicKey, R)

	cryptoText := encrypt(m[a[1]].PublicKey, originalText)
	mapToAES[a[1]] = append(mapToAES[a[1]], cryptoText)
	cryptoText = encrypt(m[a[2]].PublicKey, originalText)
	mapToAES[a[2]] = append(mapToAES[a[2]], cryptoText)
	cryptoText = encrypt(m[a[0]].PublicKey, originalText)
	mapToAES[a[0]] = append(mapToAES[a[0]], cryptoText)
	fmt.Println(cryptoText)

	text, err := decrypt(m[a[0]], mapToAES[a[1]][0])
	if (err == nil) {
		fmt.Println(text)
	} else {
		fmt.Println("Wrong private key to decode", a[0])
	}
	text, err = decrypt(m[a[1]], mapToAES[a[1]][0])
	if (err == nil) {
		fmt.Println(text)
	} else {
		fmt.Println("Wrong private key to decode", a[0])
	}
	text, err = decrypt(m[a[2]], mapToAES[a[2]][0])
	if (err == nil) {
		fmt.Println(text)
	} else {
		fmt.Println("Wrong private key to decode", a[0])
	}

	fmt.Println("number of keys in AESmap: ", len(m))
}

//construct patient struct
func newPatient(hosID, patientID string) Patient {
    return Patient{ HosID: hosID, PatientID: patientID}
}

// encrypt string to base64 crypto using AES
func encrypt(key rsa.PublicKey, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)
	label := []byte("testing")

	//keyByte, _ := x509.MarshalPKIXPublicKey(key)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key, plaintext, label)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	return string(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key *rsa.PrivateKey, cryptoText string) (string,error) {

	plaintext := []byte(cryptoText)
	label := []byte("testing")

	//keyByte, _ := x509.MarshalPKIXPublicKey(key)
	ciphertext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, plaintext, label)
	if err != nil {
		//fmt.Println("Fatal error ", err.Error())
		ciphertext = nil
		return string(ciphertext) , err
	}

	return string(ciphertext), nil
}
