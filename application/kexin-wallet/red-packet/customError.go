package red_packet

import(
	"red-envelope/network"
)

type CustomError struct {
	ErrorCode int
	//error
	ErrInfo string
	network.Failurable
}

//func (cerror *CustomError) Error() string{
//	return cerror.ErrorInfo
//}
