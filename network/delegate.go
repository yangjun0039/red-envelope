package network

type Delegate struct {
	Code         string
	Desc         string
	IsNormalREST bool // 是否正常的RESTful API模式
}

func (d Delegate) IsEqualTo(another Delegate) bool {
	return d.Code == another.Code
}
