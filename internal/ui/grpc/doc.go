// Package grpc contains all gRPC servers that expose our Application.
//
// Recommendations:
//
// 1) All files in the pkg grpc must be in a simple file list.
//
//   - order_server.go
//   - invoice_server.go
//   - user_server.go
//
// 2) If you need to convert large structures to a gRPC object and via versa, then
// create a separate file <server name>_converter.go and save this code in this
// file. Example:
//
//   - order_server.go
//   - order_converter.go
package grpc
