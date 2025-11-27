package mockrepository

import (
	"database/sql"
	"order-service/src/repository"

	"github.com/stretchr/testify/mock"
)

type TxBeginner struct {
	Mock mock.Mock
	Tx   *Transactioner
}

func NewTxBeginner() *TxBeginner {
	return &TxBeginner{
		Mock: mock.Mock{},
		Tx:   NewTransactioner(),
	}
}

func (m *TxBeginner) Begin(opts ...*sql.TxOptions) (repository.Transactioner, error) {
	args := m.Mock.Called(opts)
	return m.Tx, args.Error(1)
}
