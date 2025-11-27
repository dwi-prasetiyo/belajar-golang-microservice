package test

import (
	"context"
	"order-service/src/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestFindOrder_Repository$ -v ./src/repository/test/ -count=1

type FindOrder_TestSuite struct {
	suite.Suite
	mockDB          sqlmock.Sqlmock
	db              *gorm.DB
	orderRepository repository.Order
}

func (s *FindOrder_TestSuite) SetupTest() {
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

func (s *FindOrder_TestSuite) TearDownTest() {
	err := s.mockDB.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *FindOrder_TestSuite) Test_Success() {
	ctx := context.Background()

	whereClause := "id = ?"
	args := []any{"order-123"}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "gross_amount", "status", "payment_method", "payment_url", "created_at", "updated_at",
	}).AddRow(
		args[0],
		"user123",
		10000,
		"PAID",
		"midtrans",
		"exampleurl",
		time.Now(),
		time.Now(),
	)

	s.mockDB.ExpectQuery(`SELECT (.+) FROM "orders"`).
		WithArgs(args[0], sqlmock.AnyArg()).
		WillReturnRows(rows)

	order, err := s.orderRepository.Find(ctx, whereClause, args)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), order)
	assert.Equal(s.T(), args[0], order.ID)
}

func (s *FindOrder_TestSuite) Test_NotFound() {
	ctx := context.Background()

	whereClause := "id = ?"
	args := []any{"order-123"}

	s.mockDB.ExpectQuery(`SELECT (.+) FROM "orders"`).
		WithArgs(args[0], sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	order, err := s.orderRepository.Find(ctx, whereClause, args)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), order)
}

func (s *FindOrder_TestSuite) Test_Failed() {
	ctx := context.Background()

	whereClause := "id = ?"
	args := []any{"order-123"}

	s.mockDB.ExpectQuery(`SELECT (.+) FROM "orders"`).
		WithArgs(args[0], sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	order, err := s.orderRepository.Find(ctx, whereClause, args)
	assert.EqualError(s.T(), err, assert.AnError.Error())
	assert.Nil(s.T(), order)
}

func TestFindOrder_Repository(t *testing.T) {
	suite.Run(t, new(FindOrder_TestSuite))
}
