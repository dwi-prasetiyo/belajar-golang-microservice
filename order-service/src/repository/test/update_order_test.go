package test

import (
	"context"
	"errors"
	"order-service/src/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestUpdateOrder_Repository$ -v ./src/repository/test/ -count=1

type UpdateOrder_TestSuite struct {
	suite.Suite
	mockDB          sqlmock.Sqlmock
	db              *gorm.DB
	orderRepository repository.Order
}

func (s *UpdateOrder_TestSuite) SetupTest() {
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

func (s *UpdateOrder_TestSuite) TearDownTest() {
	err := s.mockDB.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *UpdateOrder_TestSuite) Test_Success() {
	ctx := context.Background()

	whereClause := "id = ?"
	args := []any{"order-123"}
	data := map[string]any{
		"status": "PAID",
	}

	s.mockDB.ExpectBegin()

	s.mockDB.ExpectExec(`UPDATE "orders" SET`).WithArgs(
		data["status"],
		sqlmock.AnyArg(),
		args[0],
	).WillReturnResult(sqlmock.NewResult(1, 1))

	s.mockDB.ExpectCommit()

	err := s.orderRepository.Update(ctx, whereClause, args, data)
	assert.NoError(s.T(), err)
}

func (s *UpdateOrder_TestSuite) Test_Failed() {
	ctx := context.Background()

	whereClause := "id = ?"
	args := []any{"order-123"}
	data := map[string]any{
		"status": "PAID",
	}

	s.mockDB.ExpectBegin()

	s.mockDB.ExpectExec(`UPDATE "orders" SET`).WithArgs(
		data["status"],
		sqlmock.AnyArg(),
		args[0],
	).WillReturnError(errors.New("update error"))

	s.mockDB.ExpectRollback()

	err := s.orderRepository.Update(ctx, whereClause, args, data)
	assert.EqualError(s.T(), err, "update error")
}

func TestUpdateOrder_Repository(t *testing.T) {
	suite.Run(t, new(UpdateOrder_TestSuite))
}
