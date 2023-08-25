// Package query contains APi that query to application domains.
// For example "FindRouteByID", "ActionDescription", "AvailableServices"
// etc. The query should return aggregated copy of domain object, not the domain
// objects itself.
//
// Recommendations:
// All files in the pkg query must be in a simple file list.
//
//   - find_route_by_id.go
//
//   - action_description.go
//
//   - available_services.go
//
// The query must be named the same as the file.
//
//   - action_description.go
//
//     type ActionDescription struct {}
//
//     type ActionDescriptionParameters struct {}
//
//     Each command must have the same entrypoint, func Do(). Example:
//
//     func (c *ActionDescription) Do(ctx context.Context, params ActionDescriptionParameters) error
package query
