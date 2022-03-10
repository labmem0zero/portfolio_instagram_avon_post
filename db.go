package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var db *sql.DB

func DBGetGood(goodid string)Good{
	resGood:=Good{}
	getSchema:=`
SELECT goodid, goodcategory,goodname,gooddescription,currentcat,currentcatpage,currentprice,goodurl FROM goods WHERE goodid='`+goodid+`'
`
	rows,err:=db.Query(getSchema)
	if err!=nil{
		fmt.Printf("При считывании из БД товара с picid %v произошла ошибка: %v\n",goodid,err)
	}
	for rows.Next(){
		rows.Scan(&resGood.Good.Goodid,&resGood.Good.Goodcategory,&resGood.Good.Goodname,&resGood.Good.Gooddescription,&resGood.Good.Currentcat,&resGood.Good.Currentcatpage,&resGood.Good.Currentprice,&resGood.Good.Goodurl)
	}
	return resGood
}

func DBGetCategoryGoods(category string)[]string {
	res:=[]string{}
	showSchema:=`
	SELECT goodid FROM goods WHERE goodcategory='`+category+`'
`
	rows,err:=db.Query(showSchema)
	if err!=nil{
		fmt.Println("Ошибка при считывании товаров из бд: ",err)
	}
	goodid:=""
	for rows.Next(){
		err=rows.Scan(&goodid)
		if err!=nil{
			fmt.Println("При считывании строки произошла ошибка:",err)
		}
		res=append(res,goodid)
	}
	return res
}

func DBShowByCategoryCount(){
	showSchema:=`
	SELECT goodcategory,COUNT(*) FROM goods GROUP BY goodcategory
`
	rows,err:=db.Query(showSchema)
	if err!=nil{
		fmt.Println("Ошибка при считывании товаров из бд: ",err)
	}
	goodcat:=""
	goodcount:=0
	for rows.Next(){
		err=rows.Scan(&goodcat,&goodcount)
		if err!=nil{
			fmt.Println("При считывании строки произошла ошибка:",err)
		}
		fmt.Printf("Категория: %v, количество: %v\n",goodcat,goodcount)
	}
}

func DBInsertGood(good Good){
	insertSchema:=`
	INSERT INTO goods (
		goodid, goodcategory,goodname,gooddescription,currentcat,currentcatpage,currentprice,goodurl
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?
	)
`
	_,err:=db.Exec(insertSchema, good.Good.Goodid, good.Good.Goodcategory, good.Good.Goodname, good.Good.Gooddescription, good.Good.Currentcat, good.Good.Currentcatpage,good.Good.Currentprice, good.Good.Goodurl)
	if err!=nil{
		fmt.Println("Ошибка при добавлении записи в таблицу 'goods'", err)
	}
}

func DBStart()*sql.DB{
	DBCreateGoodsSchema:=`
	CREATE TABLE IF NOT EXISTS goods (
		goodid TEXT UNIQUE,
		goodcategory TEXT,
		goodname TEXT,
		gooddescription TEXT,
		currentcat TEXT,
		currentcatpage TEXT,
		currentprice TEXT,
		goodurl TEXT
)
`

	dbPath:="C:\\avon\\goods\\goods.db"
	if _, err := os.Stat("C:\\avon\\goods\\"); os.IsNotExist(err) {
		err := os.MkdirAll("C:\\avon\\goods\\", os.ModePerm)
		if err != nil {
			fmt.Println("Ошибка при срздании папки с БД: ",err)
		}
	}
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Ошибка при создании/открытии бд: ",err)
	}
	_,err=sqlDB.Exec(DBCreateGoodsSchema)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы 'goods' в БД: ",err)
	}
	return sqlDB
}
