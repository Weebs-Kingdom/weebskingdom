package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"time"
)

func rsaConfigSetup() *rsa.PrivateKey {
	priv, err := ioutil.ReadFile("main/private.pem")
	if err != nil {
		log.Print("No RSA private key found, generating temp one")
		return generatePrivateKey(4096)
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		log.Printf("RSA private key is of the wrong type :%s", privPem.Type)
	}
	privPemBytes = privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Printf("Unable to parse RSA private key, generating a temp one :%s", err.Error())
			return generatePrivateKey(4096)
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Printf("Unable to parse RSA private key, generating a temp one : %s", err.Error())
		return generatePrivateKey(4096)
	}
	return privateKey
}

// generatePrivateKey returns a new RSA key of bits length
func generatePrivateKey(bits int) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	log.Printf("Failed to generate signing key :%s", err.Error())
	return key
}

var privateKey *rsa.PrivateKey

func InitJwt() {
	privateKey = rsaConfigSetup()
}

func GenerateLoginToken(userId primitive.ObjectID) (string, error) {
	mapp := jwt.MapClaims{
		"userId":       userId,
		"exp":          getExpiryDate(),
		"iat":          getDateNow(),
		"what are you": "looking at?",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapp)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("Failed to generate JWT :%s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func getExpiryDate() int64 {
	return time.Now().Add(time.Hour * 24 * 2).Unix()
}

func getDateNow() int64 {
	return time.Now().Unix()
}

func ParseJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey.Public(), nil
	})
	if err != nil {
		log.Printf("Failed to parse JWT :%s", err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Printf("Failed to parse JWT :%s", err.Error())
		return nil, err
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	stringByte := string(bytes)
	return stringByte, err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
