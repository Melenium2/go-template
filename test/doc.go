// Package test contains for setup database for integration tests.
// For example, you can create "fake" domain objects here for using
// it in integration tests.
//
// Example:
//  - fake_order.go
//  type FakeOrder struct {
//		ID string
//		Number string
//  }
//
//  NewFakeOrder() FakeOrder {}
//
//  - test_suite.go
//  type DBSuite struct {}
//
//  func (s *DBSuite) SetupFakeOrder(orders ...FakeOrder) error {
//  	if err := insert(orders...); err != nil {
//         ...
//  	}
//  }
package test
