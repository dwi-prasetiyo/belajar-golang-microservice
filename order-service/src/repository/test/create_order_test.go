package test

import (
	"context"
	"errors"
	"order-service/src/common/model"
	"order-service/src/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestCreateOrder_Repository$ -v ./src/repository/test/ -count=1

type CreateOrder_TestSuite struct {
	suite.Suite
	mockDB          sqlmock.Sqlmock
	db              *gorm.DB
	orderRepository repository.Order
}

func (s *CreateOrder_TestSuite) SetupTest() {
	conn, mock, err := sqlmock.New()
	s.Require().NoError(err)

	s.mockDB = mock

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 conn,
		PreferSimpleProtocol: true,
	}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)

	s.Require().NoError(err)

	s.db = db
	s.orderRepository = repository.NewOrder(db)
}

func (s *CreateOrder_TestSuite) TearDownTest() {
	err := s.mockDB.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *CreateOrder_TestSuite) Test_Success() {
	req := s.createOrder()

	s.mockDB.ExpectBegin()

	s.mockDB.ExpectExec(`INSERT INTO "orders"`).WithArgs(
		req.ID,
		req.UserID,
		req.GrossAmount,
		req.Status,
		"",
		req.PaymentURL,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	s.mockDB.ExpectCommit()

	err := s.orderRepository.CreateOrder(context.Background(), req)
	assert.NoError(s.T(), err)
}

func (s *CreateOrder_TestSuite) Test_Failed() {
	req := s.createOrder()

	s.mockDB.ExpectBegin()

	s.mockDB.ExpectExec(`INSERT INTO "orders"`).WithArgs(
		req.ID,
		req.UserID,
		req.GrossAmount,
		req.Status,
		"",
		req.PaymentURL,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnError(errors.New("insert error"))

	s.mockDB.ExpectRollback()

	err := s.orderRepository.CreateOrder(context.Background(), req)
	assert.EqualError(s.T(), err, "insert error")
}

func (s *CreateOrder_TestSuite) createOrder() *model.Order {
	orderID, _ := gonanoid.New()
	userID, _ := gonanoid.New()

	return &model.Order{
		ID:          orderID,
		UserID:      userID,
		GrossAmount: 10000,
		Status:      "PENDING_PAYMENT",
		PaymentURL:  "exampleurl",
	}
}

func TestCreateOrder_Repository(t *testing.T) {
	suite.Run(t, new(CreateOrder_TestSuite))
}
