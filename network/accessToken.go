package network

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type AccessToken struct {
	ID           string //识别ID，用户账号
	Password     string // 用户密码
	//PublicKeyPEM []byte //通信公钥
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	Password string `json:"password, omitempty"`
}

//var secretKey = []byte("1234567890qwertyuioplkjhgfdsazxcvbnm")

func (t *AccessToken) SignedString(secretKey string) (string, error) {
	//id := fmt.Sprintf("%v:%v:%v", t.Subject(), t.Identifier(), t.Validation())
	claims := AccessTokenClaims{
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: LastOfMonth(time.Now()).Unix(),
			//ExpiresAt:time.Now().Unix()+20,
			Issuer:    "zihuan",
			Id:        t.ID,
			Subject:   t.Subject(),
		},
		t.Password,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (t *AccessToken) ParseWith(signedString string, secretKey string) error {
	claims := new(AccessTokenClaims)
	err := ParseJWTWith(signedString, secretKey, claims)
	if err != nil {
		return err
	}
	t.ID = claims.Id
	t.Password = claims.Password
	return nil
}

func (t *AccessToken) VerifiableMsg(requester *Requester, r *http.Request) (string, error) {
	return StandardVerifiableMsg(requester, r)
}

func (t *AccessToken) Verify(msg string, hexSig string) error {
	return nil
}

func (t *AccessToken) Subject() string {
	return "AccessToken"
}

func (t *AccessToken) Identifier() string {
	return t.ID
}

func (t *AccessToken) Validation() string {
	return t.Password
}
