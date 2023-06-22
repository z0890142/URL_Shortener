package handler

import mysql_driver "github.com/go-sql-driver/mysql"

func SqlErrCode(err error) int {
	mysqlErr, ok := err.(*mysql_driver.MySQLError)
	if !ok {
		return 0
	}
	return int(mysqlErr.Number)
}
