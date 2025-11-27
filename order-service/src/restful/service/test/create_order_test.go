package test

import (
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	"order-service/src/common/errors"
	"order-service/src/common/model"
	"order-service/src/factory"
	mockgrpcclient "order-service/src/grpc/client/mock"
	mockrepository "order-service/src/repository/mock"
	mockrestfulclient "order-service/src/restful/client/mock"
	"order-service/src/restful/service"
	"testing"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

// go test -v ./src/restful/service/test/... -count=1 -p=1
// go test -run ^TestCreateOrder_Service$ -v ./src/restful/service/test/ -count=1

type CreateOrder_TestSuite struct {
	suite.Suite
	app             *fiber.App
	c               *fiber.Ctx
	midtransClient  *mockrestfulclient.Midtrans
	orderRepository *mockrepository.Order
	productClient   *mockgrpcclient.Product
	txRepository    *mockrepository.TxBeginner
	factory         *factory.Factory
	service         service.Order
	mocks           []*mock.Mock
}

func (s *CreateOrder_TestSuite) SetupSuite() {
	s.app = fiber.New()
	s.c = s.app.AcquireCtx(&fasthttp.RequestCtx{})

	s.midtransClient = mockrestfulclient.NewMidtrans()
	s.orderRepository = mockrepository.NewOrder()
	s.productClient = mockgrpcclient.NewProduct()
	s.txRepository = mockrepository.NewTxBeginner()

	s.mocks = []*mock.Mock{
		&s.midtransClient.Mock,
		&s.orderRepository.Mock,
		&s.productClient.Mock,
		&s.txRepository.Mock,
		&s.txRepository.Tx.Mock,
		&s.txRepository.Tx.Order.Mock,
		&s.txRepository.Tx.ProductOrder.Mock,
	}

	s.factory = &factory.Factory{
		MidtransClient:  s.midtransClient,
		OrderRepository: s.orderRepository,
		ProductClient:   s.productClient,
		TxRepository:    s.txRepository,
	}

	s.service = service.NewOrder(s.factory)
}

func (s *CreateOrder_TestSuite) TearDownSuite() {
	s.app.ReleaseCtx(s.c)
}

func (s *CreateOrder_TestSuite) TearDownTest() {
	for _, m := range s.mocks {
		m.ExpectedCalls = nil
		m.Calls = nil
	}
}

func (s *CreateOrder_TestSuite) Test_Success() {
	req := s.createOrder()

	midtransTxRes := &response.MidtransTx{
		OrderID:    "exampleid",
		PaymentURL: "exampleurl",
	}

	s.midtransClient.Mock.On(
		"CreateTransaction",
		mock.Anything,
		mock.MatchedBy(func(data *request.MidtransTransaction) bool {
			return data.TransactionDetails.OrderID != "" &&
				data.TransactionDetails.GrossAmount > 0 &&
				data.Metadata.OriginalOrderID == data.TransactionDetails.OrderID
		}),
	).Return(midtransTxRes, nil)

	s.txRepository.Mock.On(
		"Begin",
		mock.Anything,
	).Return(s.txRepository.Tx, nil)

	s.txRepository.Tx.OrderRepository().(*mockrepository.Order).On(
		"CreateOrder",
		mock.Anything,
		mock.MatchedBy(func(data *model.Order) bool {
			return data.ID != "" &&
				data.UserID == req.Order.UserID &&
				data.GrossAmount == req.Order.GrossAmount &&
				data.PaymentURL == midtransTxRes.PaymentURL &&
				data.Status == "PENDING_PAYMENT"
		}),
	).Return(nil)

	s.txRepository.Tx.ProductOrderRepository().(*mockrepository.ProductOrder).On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(data []*model.ProductOrder) bool {
			isValid := true

			for _, p := range data {
				if p.OrderID == "" || p.Quantity == 0 || p.Price == 0 {
					isValid = false
				}
			}

			return isValid
		}),
	).Return(nil)

	s.productClient.Mock.On(
		"ReduceStocks",
		mock.Anything,
		mock.MatchedBy(func(data *pb.ReduceStocksReq) bool {
			isValid := true

			for _, p := range data.ProductOrders {
				if p.ProductId == 0 || p.Quantity == 0 {
					isValid = false
				}
			}

			return isValid
		}),
	).Return(nil)

	s.txRepository.Tx.Mock.On("Commit").Return(nil)

	res, err := s.service.CreateOrder(s.c, req)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), midtransTxRes, res)
}

func (s *CreateOrder_TestSuite) Test_Failed() {
	req := s.createOrder()
	req.Order.GrossAmount = 0

	errRes := &errors.Response{HttpCode: 400, Message: "invalid request"}

	s.midtransClient.Mock.On(
		"CreateTransaction",
		mock.Anything,
		mock.MatchedBy(func(data *request.MidtransTransaction) bool {
			return data.TransactionDetails.OrderID == "" || data.TransactionDetails.GrossAmount == 0
		}),
	).Return(nil, errRes)

	res, err := s.service.CreateOrder(s.c, req)
	assert.Equal(s.T(), err.Error(), errRes.Error())

	assert.Nil(s.T(), res)
}

func (s *CreateOrder_TestSuite) createOrder() *request.CreateOrder {
	userID, _ := gonanoid.New()

	return &request.CreateOrder{
		Order: &model.Order{
			UserID:      userID,
			GrossAmount: 20000,
		},
		ProductOrders: []*model.ProductOrder{
			{
				ProductID: 101,
				Quantity:  1,
				Price:     10000,
			},
			{
				ProductID: 102,
				Quantity:  1,
				Price:     10000,
			},
		},
	}
}

func TestCreateOrder_Service(t *testing.T) {
	suite.Run(t, new(CreateOrder_TestSuite))
}
