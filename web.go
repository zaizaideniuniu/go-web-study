package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var talk string = "this is a web site,you can talk to me"

const pi = 3.1415926 //常量定义 pi

//定义数据库的配置信息
func connector() (*sql.DB, error) {
	var err error
	var db *sql.DB
	db,err = sql.Open("mysql", "root:123@/boy?charset=utf8")
	if err != nil {
		panic(err)
	}
	fmt.Println("连接Mysql数据库初始化")
	log.Println("this is my first go web service")
	return db,err
}

func check(err error) {
	if err != nil{
		fmt.Println(err)
		panic(err)
	}
}

//查询用户信息
func query() ([]byte){
	var ret []byte
	db, err := connector()
	check(err)

	rows, err := db.Query("SELECT * FROM boy.user")
	check(err)

	for rows.Next() {
		columns, _ := rows.Columns()

		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		//将数据保存到 record 字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)

		ret,_ = json.Marshal(record)
	}
	rows.Close()
	return ret
}

//插入数据 输入和输出均可参数化
/*func insertData(user *User,insert_sql string) (int64) {
	db, err := connector()
	check(err)
	stmt, err := db.Prepare(insert_sql)
	check(err)
	res, err := stmt.Exec(user.id, user.name, user.age)
	check(err)
	id, err := res.LastInsertId()
	log.Println(id)
	return id
}*/

//物理删一条数据
func remove(id string) (int64) {
	db, err := connector()
	check(err)
	stmt, err := db.Prepare("delete from user where id=?")
	check(err)
	res, err := stmt.Exec(id)
	check(err)
	num, err := res.RowsAffected()
	check(err)
	fmt.Println(num)
	stmt.Close()
	return num
}

/*func main() {
	var insert_s ="INSERT user SET id=?,name=?,age=?"
	user := newUser("114","令狐冲",27)
	insertData(user,insert_s)
	query()
	//num := remove("112")
	//log.Println("删除数据的结果是 %d",num)
}*/

func hello (w http.ResponseWriter, r * http.Request)  {
	fmt.Fprintf(w,"Hello Go web!")
}

func userHandler (w http.ResponseWriter, r * http.Request) {
	users :=query()
	fmt.Fprintf(w,"Get User LIST!")
	fmt.Fprintf(w,string(users))
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8990",
	}

	http.HandleFunc("/hello",hello)
	http.HandleFunc("/user/list",userHandler)
	server.ListenAndServe()
}
