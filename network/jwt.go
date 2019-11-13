package network

import (
	"time"
	"net/http"
	"red-envelope/foundation"
	"net/url"
	"sort"
	"github.com/pkg/errors"
	"strconv"
	"io/ioutil"
	"bytes"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

type Requester struct {
	UUID      string
	User      UserInfo
	From      Delegate
	Subject   string
	Token     Tokenable
	timestamp time.Time
}

// 封装基础日志
func (r *Requester) BeginToLog() *foundation.Entry {
	return SharedManager.Logger.WithField("uuid", r.UUID).WithField("delegate", r.From)
}

// 用户信息
type UserInfo struct {
	ID   string
	Name string
}

type Tokenable interface {
	SignedString(secretkey string) (string, error)
	ParseWith(signedString string, secretKey string) error
	VerifiableMsg(requester *Requester, r *http.Request) (string, error)
	Verify(msg string, sig string) error
	Subject() string
	Identifier() string
	Validation() string
}

func StandardVerifiableMsg(requester *Requester, r *http.Request) (string, error) {
	timestamp := r.Header.Get("X-Timestamp")
	if timestamp == "" {
		return "", errors.New("缺少参数X-Timestamp")
	}
	int64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		// 应为Unix时间戳
		return "", errors.New("X-Timestamp不合规")
	}
	remoteTime := time.Unix(int64, 0)
	duration := time.Since(remoteTime)
	if duration > 60*60*time.Second || duration < -60*60*time.Second {
		return "", errors.New("X-Timestamp已过期， " + remoteTime.Format("2006-01-02 15:04:05"))
	}
	if err := r.ParseForm(); err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	requester.BeginToLog().WithField("PostForm", r.PostForm).
		WithField("Query", r.URL.Query()).
		WithField("Body", bodyString).
		Info("VerifySignatureHandler", "request parameters")

	var id string
	if requester.Token != nil {
		id = requester.Token.Identifier()
	}
	msg := paste(
		id,
		r.Method,
		r.URL.Path,
		timestamp,
		flat(r.URL.Query()),
		flat(r.PostForm),
		bodyString,
	)
	return msg, nil
}

func paste(a ...string) string {
	sep := "|"
	return sep + strings.Join(a, sep) + sep
}

func flat(values url.Values) (str string) {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for index, key := range keys {
		str += key + "="
		str += values.Get(key)
		if index != len(keys)-1 {
			str += "&"
		}
	}
	return str
}

func StandardVerify(msg string, sig string) error {
	signedString := foundation.MD5(msg)
	if signedString != sig {
		return errors.New("signed string of msg by md5 algorithm isn't equal to signature")
	}
	return nil
}

func LastOfMonth(t time.Time) time.Time {
	currentYear, currentMonth, _ := t.Date()
	location := t.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 999999999, location)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth
}

// 解析签发的token
func ParseJWTWith(signedString string, secretKey string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}

	// token有效性验证
	if token.Valid {
		return nil
	} else {
		return fmt.Errorf("JWTSignedString[%v] fail to parse mobile claims with secret key[%v]", signedString, secretKey)
	}
}

func ParseStandardClaims(signedString string) *jwt.StandardClaims {
	claims := new(jwt.StandardClaims)
	jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error){
		return nil, nil
	})
	return claims
}
