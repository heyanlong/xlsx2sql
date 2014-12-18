package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"github.com/widuu/goini"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 获取配置
var conf = goini.SetConfig("./config.ini")

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// path
	path := conf.GetValue("path", "path")

	files, _ := filepath.Glob(path + "*")
	isInsert := make(chan bool)
	filePath := make(chan string)
	tabName := make(chan string)

	log := ""
	fmt.Println("正在导出...")

	// 忽略文件配置
	ignoreNames := strings.Split(conf.GetValue("ignore", "names"), ",")

	for _, file := range files {

		ignore := false
		for _, name := range ignoreNames {
			if file == path+name {
				ignore = true
			}
		}

		if ignore {
			continue
		}

		fmt.Println("创建一个协程处理" + file)
		go ins(file, isInsert, filePath, tabName)
	}

	for _, f := range files {
		ignore := false
		for _, name := range ignoreNames {
			if f == path+name {
				ignore = true
			}
		}

		if ignore {
			continue
		}

		isIns := <-isInsert
		file := <-filePath
		tab := <-tabName

		if isIns {
			log += file + "-->" + tab + "-->成功\r\n"
		} else {
			log += file + "-->" + tab + "-->失败\r\n"
		}
		fmt.Print(".")
	}

	ioutil.WriteFile("导出日志.log", []byte(log), 0777)
	fmt.Print("\n导出成功")
	time.Sleep(time.Second * 1)
}

func ins(file string, isInsertCh chan bool, filePathCh chan string, tabNameCh chan string) {

	xlFile, err := xlsx.OpenFile(file)

	if err != nil {
		log.Fatal(err)
	}

	isInsd := true
	tabName := ""

	for k, sheet := range xlFile.Sheets {
		if k == 0 {

			// 计算开始结束
			buffer := 500
			rows := sheet.Rows

			rowsLen := len(rows)

			remainder := rowsLen % buffer

			count := rowsLen / buffer

			if remainder > 0 {
				count += 1
			}

			tabName = getTableName(rows)
			_, data := conn()
			data.Exec("TRUNCATE TABLE " + tabName)

			for i := 0; i < count; i++ {
				start := 3
				if i > 0 {
					start = i * buffer
				}

				end := buffer

				if i > 0 {
					end = (i + 1) * buffer
				}

				if i+1 == count {
					end = rowsLen
				}

				insertSql := ""

				for k, row := range rows[start:end] {

					if k == 0 {
						insertSql = getIns(tabName)
					}

					for k3, cell := range row.Cells { // 读取列

						insertSql += "'" + strings.Trim(cell.String(), "") + "',"

						if k3 == len(row.Cells)-1 {
							insertSql = insertSql[:len(insertSql)-1]
							insertSql += "),("
						}
					}
				}

				insertSql = strings.Trim(insertSql, "")
				if len(insertSql) > 0 {
					execSql := insertSql[:len(insertSql)-2]

					_, err = data.Exec(execSql)

					if err != nil {
						isInsd = false
						break
					}
				}

				if isInsd != true {
					break
				}

			}
		}
	}

	isInsertCh <- isInsd
	filePathCh <- file
	tabNameCh <- tabName

}

func getIns(tabName string) string {

	dbname := conf.GetValue("database", "dbname")

	info, _ := conn()

	insertSql := "insert into " + tabName + " ("
	structureSql := "SELECT TABLE_NAME,COLUMN_NAME,DATA_TYPE,COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA='" + dbname + "' AND TABLE_NAME='" + tabName + "' ORDER BY TABLE_NAME,ORDINAL_POSITION"
	sqlRows, _ := info.Query(structureSql)
	for sqlRows.Next() {
		tableName, columnName, dataType, columnComment := "", "", "", ""
		sqlRows.Scan(&tableName, &columnName, &dataType, &columnComment)
		insertSql += "`" + columnName + "`,"
	}

	insertSql = insertSql[:len(insertSql)-1]

	insertSql += ") values ("
	return insertSql
}

func getTableName(rows []*xlsx.Row) string {
	for k, row := range rows {
		if k == 0 {
			for k2, tmpCell := range row.Cells {
				if k2 == 0 {
					return tmpCell.String()
				}
			}
		}
	}
	return ""
}

func conn() (info *sql.DB, data *sql.DB) {

	// DB
	username := conf.GetValue("database", "username")
	password := conf.GetValue("database", "password")
	dbname := conf.GetValue("database", "dbname")

	info, infoErr := sql.Open("mysql", username+":"+password+"@/information_schema")
	data, dataErr := sql.Open("mysql", username+":"+password+"@/"+dbname)

	if infoErr != nil || dataErr != nil {
		log.Fatal("数据库连接失败！")
	}

	return info, data
}
