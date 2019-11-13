package foundation

import
(
	"testing"
	"time"
	"fmt"
)


func TestTime(t *testing.T){
	tt := time.Now().Unix()
	fmt.Println(tt)
}