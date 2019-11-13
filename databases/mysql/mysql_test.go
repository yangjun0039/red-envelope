package mysql

import (
	//"database/sql"
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"red-envelope/configer"
	"fmt"
)

func TestFetchdata(t *testing.T) {
	sql := "select * from role"

	configer.InitConfiger()
	New()
	maps,err := FetchRowsD(sql)

	fmt.Println("*******************************************")
	fmt.Println("*******************************************")
	fmt.Println("*******************************************")

	if err != nil{
		fmt.Println(err)
	}
	for _,m := range(maps){
		fmt.Println(m)
	}
}


//func fetchData(rows *sql.Rows) ([]map[string]string, error) {
//	columns, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println(columns)
//	values := make([]sql.RawBytes, len(columns))
//	scanArgs := make([]interface{}, len(values))
//
//	ret := make([]map[string]string, 0)
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	for rows.Next() {
//		err = rows.Scan(scanArgs...)
//		if err != nil {
//			return nil, err
//		}
//		var value string
//		vmap := make(map[string]string, len(scanArgs))
//		for i, col := range values {
//			if col == nil {
//				value = ""
//			} else {
//				value = string(col)
//			}
//			vmap[columns[i]] = value
//		}
//		ret = append(ret, vmap)
//	}
//	return ret, nil
//}
