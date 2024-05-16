package tests

import (
	"context"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"password_keeper/config/server"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/server/handler"
	"password_keeper/internal/server/repository"
	"password_keeper/internal/server/service"
)

type APITestSuite struct {
	suite.Suite

	db      *sqlx.DB
	handler *handler.Handler
	serv    *service.Service
	rep     *repository.Repository
	cfg     *server.ServConfig
	enc     *encryption.Decryptor

	mocks *mocks
}

type mocks struct {
	serviceAuthMock *service.MockAuthorizationService
	serviceDataMock *service.MockDataServiceInterface
	repAuthMock     *repository.MockAuthorizationRepository
	repDataMock     *repository.MockDataRepositoryInterface
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := server.NewServer()

	s.cfg = cfg

	db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		s.FailNow("Failed to initialize database", err)
	}
	s.db = db
	s.initDeps()
	s.initMocks()
	s.testData()
}

func (s *APITestSuite) initDeps() {
	rep := repository.NewRepository(s.db)
	serv := service.NewService(rep, s.cfg)
	hand := handler.NewHandler(serv)
	err := encryption.InitDecryptor(s.cfg.PrivateCryptoKeyPath)
	if err != nil {
		s.FailNow("Failed to initialize decryptor", err)
	}

	s.serv = serv
	s.handler = hand
	s.rep = rep
}

func (s *APITestSuite) initMocks() {
	c := gomock.NewController(s.T())
	s.mocks = &mocks{
		serviceAuthMock: service.NewMockAuthorizationService(c),
		serviceDataMock: service.NewMockDataServiceInterface(c),
		repAuthMock:     repository.NewMockAuthorizationRepository(c),
		repDataMock:     repository.NewMockDataRepositoryInterface(c),
	}
}

func (s *APITestSuite) TearDownSuite() {
	s.db.Close()
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) testData() {
	ctx := context.Background()
	user := entity.User{Login: "testSignIn", Password: "testSignIn"}
	_, err := s.serv.CreateUser(ctx, user)
	if err != nil {
		s.FailNow("Failed to set user", err)
	}

}
