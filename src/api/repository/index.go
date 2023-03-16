package repository

import (
	httperors "github.com/myrachanto/erroring"
	model "github.com/myrachanto/sqlgostructure/src/api/models"
	"github.com/spf13/viper"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// IndexRepo
var (
	IndexRepo indexRepo = indexRepo{}
)

// Layout ...
const (
	Layout   = "2006-01-02"
	layoutUS = "January 2, 2006"
)

type Db struct {
	DbType     string `mapstructure:"DbType"`
	DbName     string `mapstructure:"DbName"`
	DbUsername string `mapstructure:"DbUsername"`
	DbPassword string `mapstructure:"DbPassword"`
	DbPort     string `mapstructure:"DbPort"`
}

func LoaddbConfig() (db Db, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&db)
	return
}

// /curtesy to gorm
type indexRepo struct {
	Bizname string `json:"bizname,omitempty"`
}

func (indexRepo indexRepo) Dbsetup() httperors.HttpErr {
	db := "simple"
	GormDB, err1 := gorm.Open(sqlite.Open(db), &gorm.Config{})
	if err1 != nil {
		httperors.NewNotFoundError("Failed to initialize the system!")
	}
	GormDB.AutoMigrate(&model.User{})

	return nil
}
func (indexRepo indexRepo) Getconnected() (*gorm.DB, httperors.HttpErr) {
	db := "simple"
	GormDB, err1 := gorm.Open(sqlite.Open(db), &gorm.Config{})
	if err1 != nil {
		httperors.NewNotFoundError("Failed to initialize the system!")
	}
	return GormDB, nil
}
func (indexRepo indexRepo) DbClose(GormDB *gorm.DB) {
	sqlDB, err := GormDB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
