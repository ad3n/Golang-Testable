package grpcs

import (
	"context"

	"github.com/ad3n/golang-testable/protos"
	"github.com/ad3n/golang-testable/services"
)

type Account struct {
	Service *services.Account
}

func (s *Account) Balance(c context.Context, request *protos.BalanceRequest) (*protos.BalanceResponse, error) {
	result, err := s.Service.Balance(int(request.AccountNumber))
	if err != nil {
		return &protos.BalanceResponse{}, err
	}

	response := protos.BalanceResponse{}
	response.AccountNumber = int32(result.AccountNumber)
	response.CustomerName = result.CustomerName
	response.Balance = float32(result.Balance)

	return &response, nil
}
