package network

import (
	"net/http"
	"time"
	"encoding/json"
)

type ResponseRecorder struct {
	Status    int
	uuid      string
	jsonBytes []byte
	w         http.ResponseWriter
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{Status: 200, jsonBytes: []byte{}, w: w}
}

func (self *ResponseRecorder) Header() http.Header {
	return self.w.Header()
}

func (self *ResponseRecorder) Write(b []byte) (int, error) {
	return self.w.Write(b)
}

func (self *ResponseRecorder) WriteHeader(status int) {
	self.w.WriteHeader(status)
}

func (self *ResponseRecorder) response() {
	self.WriteHeader(self.Status)
	self.Write(self.jsonBytes)
	if self.Status >= 100 && self.Status < 199 {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Info("response", "Infomstional 1xx")
	} else if self.Status >= 200 && self.Status < 299 {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Info("response", "success 2xx")
	} else if self.Status >= 300 && self.Status < 399 {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Info("response", "Redirection 3xx")
	} else if self.Status >= 400 && self.Status < 499 {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Error("response", "Client Error 4xx")
	} else if self.Status >= 500 && self.Status < 599 {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Error("response", "Server Error 5xx")
	} else {
		SharedManager.Logger.WithField("uuid", self.uuid).
			WithField("HTTPStatus", self.Status).
			WithField("response", string(self.jsonBytes)).Warn("response", "Unknow HTTP Status")
	}
}

func (self *ResponseRecorder) Reset() {
	self.Status = 200
	self.jsonBytes = []byte{}
}

type Success struct {
	status      int
	infoPointer interface{}
}

func NewSuccess(status int, infoPointer interface{}) *Success {
	return &Success{status, infoPointer}
}

func (self *Success) responseContext(w http.ResponseWriter) *ResponseRecorder {
	var bytes []byte
	var err error
	if self.infoPointer != nil {
		bytes, err = json.Marshal(self.infoPointer)
		if err != nil {
			return NewFailure(http.StatusInternalServerError, jsonSerializationFailure).responseContext(w)
		}
	}
	if recorder, ok := w.(*ResponseRecorder); ok {
		if len(recorder.jsonBytes) != 0 {
			//  日志记录 todo
			recorder.jsonBytes = []byte{}
			return recorder
		}
		recorder.Status = self.status
		recorder.jsonBytes = bytes
		return recorder
	} else {
		recorder = NewResponseRecorder(w)
		recorder.Status = self.status
		recorder.jsonBytes = bytes
		return recorder
	}

}

func (self *Success) Response(w http.ResponseWriter) {
	self.responseContext(w).response()
}

// RESTful API 默认失败报文
type Failure struct {
	status     int    `json:"-"`           // HTTP状态码
	Code       string `json:"errorCode"`   // 内部错误码
	ErrorMsg   string `json:"error_msg"`   // 错误信息
	DisplayMsg string `json:"display_msg"` // 待展示的信息
	Timestamp  string `json:"timestamp"`   // 时间戳
}

type Failurable interface {
	Code() string
	ErrorMsg() string
	DisplayedMsg() string
}

func NewFailure(status int, f Failurable) *Failure {
	return &Failure{
		status:     status,
		Code:       f.Code(),
		ErrorMsg:   f.ErrorMsg(),
		DisplayMsg: f.DisplayedMsg(),
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (self *Failure) AppendErrorMsg(str string) *Failure {
	self.ErrorMsg += ", " + str
	return self
}

func (self *Failure) SetDisplayedMsg(str string) *Failure {
	self.DisplayMsg = str
	return self
}

func (self *Failure) responseContext(w http.ResponseWriter) *ResponseRecorder {
	bytes, _ := json.Marshal(self)
	if recorder, ok := w.(*ResponseRecorder); ok {
		if len(recorder.jsonBytes) != 0 {
			//todo
		}
		recorder.Status = self.status
		recorder.jsonBytes = bytes
		return recorder
	} else {
		recorder = NewResponseRecorder(w)
		recorder.Status = self.status
		recorder.jsonBytes = bytes
		return recorder
	}
}

func (self *Failure) Response(w http.ResponseWriter) {
	self.responseContext(w).response()
}
