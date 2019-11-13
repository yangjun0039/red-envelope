package network

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strings"
	"time"
	"github.com/satori/go.uuid"
	"os"
	"runtime/debug"
	"red-envelope/application/system/model"
	"regexp"
)

type Handler func(requester *Requester, w http.ResponseWriter, r *http.Request)

// CrossOriginHandler is the middleware for allowing cross origin request.
// The white origins are set by the white_origins list in config file.
func CrossOriginHanddler(handler Handler) Handler {
	fmt.Println("CrossOriginHanddler")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		origin := r.Header.Get("Origin")
		requester.BeginToLog().WithField("origin", origin).Info("CrossOriginHanddler", "a request crossed origin")
		fmt.Println(SharedManager.whiteOrigins)
		if SharedManager.whiteOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE")
			w.Header().Set("Access-Control-Allow_heads", "Authorization, X-Method, X-Timestamp, X-Signature, Content-Type")
		}
		handler(requester, w, r)
	}
}

func RecoveryHandler(handler Handler) Handler {
	fmt.Println("RecoveryHandler")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				SharedManager.Logger.WithField("PID", os.Getpid()).
					WithField("error", err).
					WithField("stack", string(debug.Stack())).
					Error("RecoveryHandler", "a recoverable panic happened")
				errMsg := fmt.Sprintf("Panic: %v", err)
				NewFailure(http.StatusInternalServerError, unknownPanic).AppendErrorMsg(errMsg).Response(w)
			}
		}()
		handler(requester, w, r)
	}
}

// VerifySignatureHandler is the middleware for verifying signature,
// only digested verifiable msg is matched signature, the request is legal.
func VerifySignatureHandler(handler Handler) Handler {
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Signature")
		if signature == "" {
			NewFailure(http.StatusBadRequest, signatureVerificationIsFailed).AppendErrorMsg("缺少参数X-Signature").Response(w)
			return
		}
		signature = strings.ToLower(signature)
		msg, err := StandardVerifiableMsg(requester, r)
		if err != nil {
			NewFailure(http.StatusInternalServerError, signatureVerificationIsFailed).AppendErrorMsg(err.Error()).Response(w)
			return
		}
		err = StandardVerify(msg, signature)
		if err != nil {
			requester.BeginToLog().WithField("msg", msg).
				WithField("signature", signature).
				WithField("error", err).Error("VerifySignatureHandler", "fail to verify")
			NewFailure(http.StatusInternalServerError, signatureVerificationIsFailed).AppendErrorMsg(err.Error()).Response(w)
			return
		}
		handler(requester, w, r)
	}
}

func RequesterHandler(handler Handler) http.Handler {
	fmt.Println("RequesterHandler")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := uuid.NewV4()
		uuid := u.String()
		aRequester := &Requester{
			uuid,
			UserInfo{},
			Delegate{},
			"",
			nil,
			time.Now(),
		}
		aRequester.BeginToLog().WithField("url", r.URL.Path).Infof("RequesterIn", "requester %v is coming", aRequester.UUID)
		recorder := NewResponseRecorder(w)
		recorder.uuid = uuid
		handler(aRequester, recorder, r)
		duration := time.Since(aRequester.timestamp)
		aRequester.BeginToLog().WithField("duration", duration.Seconds()).Info("RequesterOut", "requester %v is leaving", aRequester.UUID)
	})
}

// 用于重定向
func SwitchRouteHandler(handler http.Handler) http.Handler {
	fmt.Println("SwitchRouteHandler")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Header.Get("x-Method")
		if method != "" && (r.Method == http.MethodGet || r.Method == http.MethodPost) {
			route := mux.CurrentRoute(r)
			name := route.GetName()
			name = strings.Replace(name, r.Method, method, 1)
			h := rootRouter.Multiplexer.GetRoute(name).GetHandler()
			r.Method = method
			r.Header.Del("X-Method")
			h.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

//func AuthorizeHandler(handler Handler) Handler {
//	fmt.Println("AuthorizeHandler")
//	return Handler(func(requester *Requester, w http.ResponseWriter, r *http.Request) {
//
//	})
//}

func ForbidRequestTooOftenHandler(handler Handler) Handler {
	fmt.Println("ForbidRequestTooOftenHandler")
	return Handler(func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		now := time.Now().UnixNano()
		diff := now
		var key string
		defer func() {
			if diff <= 1*1e9 {
				requester.BeginToLog().WithField("now", now).
					WithField("diff", diff).
					Error("ForbidRequestTooOftenHandler", "request is too often")
				NewFailure(http.StatusForbidden, requestIsTooOften).Response(w)
				return
			} else {
				err := SharedManager.RequestRecorder.SetLastTimestamp(key, now)
				if err != nil {
					requester.BeginToLog().WithField("now", now).
						WithField("key", key).
						WithField("error", err).
						Error("ForbidRequestTooOftenHandler", "request is too often")
					return
				}
				handler(requester, w, r)
			}
		}()
		if SharedManager.RequestRecorder == nil {
			fmt.Println("RequestRecorder is nil")
			return
		}
		addr := r.RemoteAddr
		uri := fmt.Sprintf("%v:%v", r.Method, r.URL.Path)
		id := "guest"

		if requester.Token != nil {
			id = requester.Token.Identifier()
		}
		key = fmt.Sprintf("%v:%v:%v", id, addr, uri)
		lastTimestamp, err := SharedManager.RequestRecorder.GetLastTimestamp(key)
		if err != nil {
			// 记录日志
			requester.BeginToLog().WithField("now", now).
				WithField("lastTimestamp", lastTimestamp).
				WithField("key", key).
				WithField("error", err).
				Warn("ForbidRequestTooOftenHandler", "no last timestamp")
			return
		}
		diff = now - lastTimestamp
	})
}

func LoginTokenHandler(handler Handler) Handler {
	fmt.Println("LoginTokenHandler")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		t := requester.Token
		accToken, ok := t.(*AccessToken)
		if !ok {
			acc := r.PostFormValue("account")
			pwd := r.PostFormValue("password")
			if acc == "" {
				NewFailure(http.StatusForbidden, invalidAccount).AppendErrorMsg("account is nil").Response(w)
				return
			}
			if pwd == "" {
				NewFailure(http.StatusForbidden, invalidPassword).AppendErrorMsg("password is nil").Response(w)
				return
			}
			requester.Token = &AccessToken{
				ID:       acc,
				Password: pwd,
			}
		} else {
			requester.Token = accToken
		}
		handler(requester, w, r)
	}
}

func JWTHandler(handler Handler) Handler {
	fmt.Println("JWTHandler")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			NewFailure(http.StatusBadRequest, lackParameter).AppendErrorMsg("缺少信息： Authorization").Response(w)
			return
		}
		claims := ParseStandardClaims(authorization)
		key := fmt.Sprintf("%v:%v", claims.Subject, claims.Id)
		fun := SharedManager.GetTokenGenerator(claims.Subject)
		if fun == nil{
			NewFailure(http.StatusUnauthorized, invalidJSONWebToken).AppendErrorMsg(claims.Subject).Response(w)
			return
		}
		t := fun()
		err := t.ParseWith(authorization, key)
		if err != nil {
			NewFailure(http.StatusUnauthorized, invalidJSONWebToken).AppendErrorMsg(err.Error()).Response(w)
			return
		}

		requester.Token = t
		handler(requester, w, r)
	}
}

// 校验权限
func VerfyPermissionsHandler(handler Handler) Handler {
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		permission, err := model.QryUserPermission(requester.Token.Identifier())
		if err != nil {
			NewFailure(http.StatusInternalServerError, invalidAccount).AppendErrorMsg(err.Error()).Response(w)
		}
		path := r.URL.Path
		permissionAllowed := checkPermission(path, r.Method, permission)
		if !permissionAllowed {
			NewFailure(http.StatusForbidden, permissionNotAllow).Response(w)
		}
		handler(requester, w, r)
	}
}

func checkPermission(path, method string, permission []map[string]string) bool {
	for _, p := range permission {
		methods := strings.Split(p["http-method"], ",")
		if p["http_method"] == "" || inMethodArr(methods, method) {
			if p["http_path"] == "*" {
				return true
			}
			if path == p["http_path"] {
				return true
			}
			reg, err := regexp.Compile(p["http_path"])
			if err != nil {
				fmt.Println("err:", err)
				continue
			}
			if reg.FindString(path) == p["http_path"] {
				return true
			}

		}
	}
	return false
}

func inMethodArr(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], str) {
			return true
		}
	}
	return false
}
