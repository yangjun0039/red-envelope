package mysql

import (
	"database/sql"
	"sync"
	"red-envelope/configer"
	"fmt"
	"time"
	"github.com/pkg/errors"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var mu sync.Mutex

var once sync.Once

var dbNilError = errors.New("db is nil")

func GetDB() *sql.DB {
	if db == nil {
		openDB()
	}
	return db
}

func New() {
	if db == nil {
		once.Do(openDB)
	} else {
		fmt.Println("db already exists")
	}
}

func openDB() {
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		return
	}
	var err error
	db, err = sql.Open("mysql", configer.MySqlConfig().DataSourceName)
	if err != nil {
		errMsg := fmt.Sprintf("数据库连接失败1：%v", err)
		fmt.Println(errMsg)
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * 10)
	if err = db.Ping(); err != nil {
		errMsg := fmt.Sprintf("数据库连接失败2：%v", err)
		fmt.Println(errMsg)
		panic(err)
	}

}

// 插入
func Insert(sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	stmtIns, err := db.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 插入(有事务处理) 调用时不用Rollback
func InsertTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

// 修改和删除(有事务处理) 调用时不用Rollback
func ExecTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

// 修改和删除
func Exec(db *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	stmtIns, err := db.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 查询一行数据
func FetchRowD(sqlstr string, args ...interface{}) (map[string]string, error) {
	if db == nil {
		return nil, dbNilError
	}
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return map[string]string{}, nil
}

// 查询多行数据
func FetchRowsD(sqlstr string, args ...interface{}) ([]map[string]string, error) {
	if db == nil {
		return nil, dbNilError
	}
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func fetchData(rows *sql.Rows) ([]map[string]string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return ret, nil
}
