package peer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type FileClaim struct {
	//Size of file in bytes
	Size int64 `json:"size"`

	//Filename
	Name string `json:"name"`

	//Action - 0=read; 1=write; 2=delete
	Act int8 `json:"act"`

	jwt.StandardClaims
}

const (
	READACT = 0
	WRITACT = 1
	DELEACT = 2
)

type FileValidationResponse struct {
	//Validation status: 0=valid, 1=invalid
	Valid bool `json:"status"`

	//Validation message
	Message string `json:message`
}

func getBaseName(fname string) (string, error) {
	pattern, err := regexp.Compile("((_[[:lower:]]*)?(_rep[[:digit:]]+))?$")
	if err != nil {
		return "", err
	}

	indices := pattern.FindAllStringIndex(fname, -1)
	if indices == nil {
		return fname, nil
	}

	suff_begin := indices[len(indices)-1][0]
	return fname[:suff_begin], nil
}

// DecodeCertificate decodes the tokenString without validating the signature
func decodeCertificate(tokenString string) (int64, string, int8, error) {
	var claims_class jwt.Claims = &FileClaim{}
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, _, err := parser.ParseUnverified(tokenString, claims_class)

	if err != nil {
		return -1, "", -1, err
	}

	claims := token.Claims.(*FileClaim)

	return claims.Size, claims.Name, claims.Act, nil
}

func validateCertificate(tokenString string) error {

	responseRaw, err := http.Post("http://172.104.136.183/auth/node/action", "text/plain", strings.NewReader(tokenString))
	if err != nil {
		return err
	}

	responseParsed := FileValidationResponse{}
	jsonDec := json.NewDecoder(responseRaw.Body)
	jsonDec.Decode(&responseParsed)

	if !responseParsed.Valid {
		return fmt.Errorf("Certificate invalidated by server, msg=%s", responseParsed.Message)
	}

	return nil
}

func ValidateFile(shardname string, tokenString string, action int8) error {

	basename, err := getBaseName(shardname)
	if err != nil {
		return err
	}

	fsize_cert, basename_cert, action_cert, err := decodeCertificate(tokenString)
	if err != nil {
		return err
	}

	if action_cert != action {
		return fmt.Errorf("Actions in certificate and request don't match: %d != %d", action_cert, action)
	}
	if basename_cert != basename {
		return fmt.Errorf("Certificate name doesn't match request name: %s != %s", basename_cert, basename)
	}

	// No need to check filesize when writing (server does it for us)
	if action == WRITACT {
		return nil
	}

	// Check file size
	fi, err := os.Stat(shardname)
	if err != nil {
		return err
	}

	fsize := fi.Size()
	if fsize > fsize_cert {
		return fmt.Errorf("Certificate file size doesn't match: %d != %d", fsize_cert, fsize)
	}

	if err := validateCertificate(tokenString); err != nil {
		return err
	}

	return nil
}
