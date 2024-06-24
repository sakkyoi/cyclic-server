package magistrate

import (
	"crypto/rsa"
	"cyclic/pkg/colonel"
	"cyclic/pkg/scribe"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"go.uber.org/zap"
	"net/http"
	"slices"
	"time"
)

type Magistrate struct {
	signKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
}

type Claims struct {
	jwt.RegisteredClaims
	Info
}

type Info struct {
	ID string `json:"id"`
}

func New() *Magistrate {
	m := &Magistrate{}
	m.signKey, m.VerifyKey = m.init()
	return m
}

func (m *Magistrate) init() (*rsa.PrivateKey, *rsa.PublicKey) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(colonel.Writ.JWT.PrivateKey))
	if err != nil {
		scribe.Scribe.Fatal("failed to parse private key", zap.Error(err))
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(colonel.Writ.JWT.PublicKey))
	if err != nil {
		scribe.Scribe.Fatal("failed to parse public key", zap.Error(err))
	}

	return signKey, verifyKey
}

func (m *Magistrate) Issue(aud []string, id string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &Claims{
		jwt.RegisteredClaims{
			Issuer:    "cyclic",
			Subject:   id,
			Audience:  aud,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(colonel.Writ.JWT.Expiration))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        time.Now().String(), // TODO: implement jwt id
		},
		Info{
			ID: id,
		},
	}

	return t.SignedString(m.signKey)
}

func (m *Magistrate) Gavel(r *http.Request) (*Claims, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor,
		func(token *jwt.Token) (interface{}, error) {
			return m.VerifyKey, nil
		}, request.WithClaims(&Claims{}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}

func (m *Magistrate) Examine(claims *Claims, aud string) bool {
	if !slices.Contains(claims.Audience, aud) {
		return false
	}

	return true
}
