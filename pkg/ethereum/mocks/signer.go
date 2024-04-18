package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

type Signer struct {
	mock.Mock
}

func (s *Signer) Address() ethereum.Address {
	args := s.Called()
	return args.Get(0).(ethereum.Address)
}

func (s *Signer) SignTransaction(transaction *ethereum.Transaction) error {
	args := s.Called(transaction)
	return args.Error(0)
}

func (s *Signer) Signature(data []byte) (ethereum.Signature, error) {
	args := s.Called(data)
	return args.Get(0).(ethereum.Signature), args.Error(1)
}

func (s *Signer) Recover(signature ethereum.Signature, data []byte) (*ethereum.Address, error) {
	args := s.Called(signature, data)
	return args.Get(0).(*ethereum.Address), args.Error(1)
}
