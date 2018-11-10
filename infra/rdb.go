package infra

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	// Just import.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// RDB provides gorm db reference.
	RDB *gorm.DB
)

type dbConfig struct {
	Endpoint string
	Database string
	Username string
	Password string
}

var dbConfigs = map[environment]dbConfig{
	dev: dbConfig{
		Database: "slgmails_dev",
	},
	stg: dbConfig{
		Database: "slgmails_stg",
	},
	prod: dbConfig{
		Database: "slgmails_prod",
	},
}

func setupDatabase() {
	// Determine base config per env.
	dbc := dbConfigs[getEnvironment()]

	Linfo("Get RDS username & password from env...")
	dbc.Username = Env.MysqlUser
	dbc.Password = Env.MysqlPass

	Linfo("Get MYSQL endpoint...")

	dbc.Endpoint = Env.MysqlEndpoint

	if dbc.Endpoint == "" {
		Lpanic("db endpoint not found on environment variable and config")
	}

	Linfo("Connect database...")
	connectToDatabase(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			dbc.Username, dbc.Password, dbc.Endpoint, dbc.Database),
	)
}

func connectToDatabase(dns string) {
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		Lpanic(fmt.Sprintf("%s", err))
	}

	db.LogMode(!IsProduction())

	// SetMaxIdleConns() should be greater than SetMaxOpenConns().
	db.DB().SetMaxIdleConns(30)
	db.DB().SetMaxOpenConns(30)
	// SetConnMaxLifetime() should be MaxOpenConns * 1sec.
	db.DB().SetConnMaxLifetime(time.Second * 30)
	RDB = db
}
