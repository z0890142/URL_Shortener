package common

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

func OpenMysqlDatabase(option *config.DatabaseOption) (db *sql.DB, err error) {

	connection, err := GetConnectionString(option)
	if err != nil {
		return nil, fmt.Errorf("openMysqlDatabase: %v", err)
	}

	if db, err = sql.Open(c.DriverMysql, connection); err != nil {
		return nil, fmt.Errorf("openMysqlDatabase: %v", err)
	} else {
		err = db.Ping()
		if err != nil {
			return nil, fmt.Errorf("openMysqlDatabase: %v", err)
		}
	}

	// Set connection pool
	if option.PoolSize > 0 {
		db.SetMaxIdleConns(option.PoolSize)
		db.SetMaxOpenConns(option.PoolSize)
	}

	return
}

func GetConnectionString(option *config.DatabaseOption) (string, error) {

	var loc = time.Local
	var err error
	if len(option.Timezone) > 0 {
		if loc, err = time.LoadLocation(option.Timezone); err != nil {
			return "", fmt.Errorf("GetConnectionString: %v", err)
		}
	}
	c := mysql.Config{
		User:                 option.Username,
		Passwd:               option.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", option.Host, option.Port),
		DBName:               option.DBName,
		Loc:                  loc,
		Timeout:              option.Timeout,
		ReadTimeout:          option.ReadTimeout,
		WriteTimeout:         option.WriteTimeout,
		ParseTime:            true,
		CheckConnLiveness:    true,
		AllowNativePasswords: true,
		MaxAllowedPacket:     4 << 20, // 4MB
		Collation:            "utf8mb4_general_ci",
		MultiStatements:      true,
	}
	if len(option.Charset) > 0 {
		c.Params = make(map[string]string)
		c.Params["charset"] = option.Charset
	}
	return c.FormatDSN(), nil

}
