// Package container contains source code for initializing all dependencies of
// out application. A bit more for each file included in this pkg.
//
// File config.go
// Config contains structs that will be filled during application startup.
// At the current moment, config get all the variables from environment
// variables. List of environment variable for the local development
// can be specified in .env file. Also, environment variables can be
// in OS environment variables, then they will also appear in the
// Config struct.
//
// File container.go
// Container initialize all the application dependencies and stores
// them in defined structures. The Container contains some default
// structures for storing dependencies inside, but you can define
// your own.
//
// File dependencies.go
// Contains definition of each application dependency. All the functions
// use the Container for setup own dependencies.
package container
