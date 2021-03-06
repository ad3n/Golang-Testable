// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/ad3n/golang-testable/models"

	mock "github.com/stretchr/testify/mock"
)

// CustomerRepository is an autogenerated mock type for the CustomerRepository type
type CustomerRepository struct {
	mock.Mock
}

// Find provides a mock function with given fields: Id
func (_m *CustomerRepository) Find(Id int) (*models.Customer, error) {
	ret := _m.Called(Id)

	var r0 *models.Customer
	if rf, ok := ret.Get(0).(func(int) *models.Customer); ok {
		r0 = rf(Id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Customer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(Id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Saves provides a mock function with given fields: customers
func (_m *CustomerRepository) Saves(customers ...*models.Customer) error {
	_va := make([]interface{}, len(customers))
	for _i := range customers {
		_va[_i] = customers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*models.Customer) error); ok {
		r0 = rf(customers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
