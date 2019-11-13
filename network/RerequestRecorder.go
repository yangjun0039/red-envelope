package network

import
(
	"red-envelope/databases/redis"
)

type ReRequestRecorder struct {

}

func (re ReRequestRecorder)SetLastTimestamp(key string, value int64) error{
	_, err := redis.GetInstance().Set(key, value, 60)
	return err
}
func (re ReRequestRecorder)GetLastTimestamp(key string) (int64, error){
	return redis.GetInstance().GetInt64(key)
}
