// Package command contains APi that combine different application domains
// into one use case. For example "CreateOrder", "MakeInvoice", "CreateClient"
// etc. If possible, the command should be idempotent and return no value other
// than an error.
//
// Recommendations:
// All files in the pkg command must be in a simple file list.
//   - create_order.go
//   - make_invoice.go
//   - create_client.go
//
// The command must be named the same as the file.
//
//   - create_order.go
//
//     type CreateOrder struct {}
//
//     type CreateOrderParameters struct {}
//
// Each command must have the same entrypoint, func Do(). Example:
//
//	func (c *CreateOrder) Do(ctx context.Context, params CreateOrderParameters) error
package command
