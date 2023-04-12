package storage

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"time"
)

type Service struct {
	Db     *gorm.DB
	Config *Config
	logger *zap.SugaredLogger
}

type Config struct {
	Db              *Db
	ScanApiEndpoint string
	ScanApiKey      string
	RpcWsEndpoint   string
}

type Db struct {
	User     string
	Password string
	Database string
	Host     string
	Port     int
}

func NewService(logger *zap.SugaredLogger) *Service {
	config := &Config{
		Db: &Db{
			User:     "bot",
			Password: "password",
			Database: "bot",
			Host:     "localhost",
			Port:     5432,
		},
		ScanApiEndpoint: "",
		ScanApiKey:      "",
		RpcWsEndpoint:   "",
	}

	return &Service{
		Config: config,
		logger: logger,
	}
}

func (s *Service) Start() {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		s.Config.Db.Host,
		s.Config.Db.User,
		s.Config.Db.Password,
		s.Config.Db.Database,
		s.Config.Db.Port)

	dataBase, err := gorm.Open(
		postgres.Open(DSN),
		&gorm.Config{Logger: zapgorm2.New(s.logger.Desugar())},
	)
	if err != nil {
		s.logger.With(err).Fatal("Can't connect to DB")
	}

	sqlDB, err := dataBase.DB()
	if err != nil {
		s.logger.With(err).Info("DB connection error")
	} else if sqlDB != nil {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour * 240)
	} else {
		s.logger.Info("DB reconnect existsDB is empty")
	}

	s.Db = dataBase
	s.logger.Info("DB connected")
}
