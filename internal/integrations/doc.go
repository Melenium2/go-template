// Package integrations contains source code for integration
// with other systems that not included in our list of application.
// For example, some system for sending SMS, routing system etc.
//
// This package is anti-corruption layer for business domains of our
// application. This package should contain the code that converts external
// entities into our entities. For example, "Order" from the  Maxoptra
// system and "Order" from our system.
package integrations
