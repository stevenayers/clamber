package service_test

import (
	"github.com/stevenayers/clamber/service"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"strings"
	"testing"
)

type (
	StoreSuite struct {
		suite.Suite
		store   service.DbStore
		crawler service.Crawler
		logger  log.Logger
		config  service.Config
	}
)

func (s *StoreSuite) SetupSuite() {
	var err error
	configFile := "../test/Config.toml"
	s.config, err = service.InitConfig(configFile)
	if err != nil {
		s.T().Fatal(err)
	}

	s.logger = InitJsonLogger(log.NewSyncWriter(os.Stdout), s.config.General.LogLevel)
	s.store = service.DbStore{}
	s.store.Connect(s.config.Database)
}

func (s *StoreSuite) SetupTest() {
	var err error
	configFile := "../test/Config.toml"
	s.config, err = service.InitConfig(configFile)
	if err != nil {
		s.T().Fatal(err)
	}
	s.store.Connect(s.config.Database)
	if !strings.Contains(s.T().Name(), "TestLog") && !strings.Contains(s.T().Name(), "TestConnect") {
		err := s.store.DeleteAll()
		if err != nil {
			s.T().Fatal(err)
		}
		err = s.store.SetSchema()
		if err != nil {
			s.T().Fatal(err)
		}
	}
}

func (s *StoreSuite) TearDownSuite() {
	for _, conn := range s.store.Connection {
		err := conn.Close()
		if err != nil {
			fmt.Print(err)
		}
	}
}

func TestSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func InitJsonLogger(writer io.Writer, logLevel string) (logger log.Logger) {
	logger = log.NewJSONLogger(writer)
	logger = log.With(
		logger,
		"service", "clamber-api",
		"node", uuid.New().String(),
	)
	switch logLevel {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	return
}
