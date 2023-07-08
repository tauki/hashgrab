// Package hashgrab provides utilities to fetch data from URLs,
// hash the data using various algorithms, and manage the process
// with a configurable number of concurrent operations.
//
// This package includes interfaces for fetchers and hashers,
// allowing users to define their own custom implementations
// for these components. It also includes a Worker type
// that coordinates fetching and hashing operations.
//
// The typical use case is to create a new Worker using the New function,
// configure the Worker as needed using its methods,
// and then call the Worker's Run method with a list of URLs to process.
// Results of the operations are returned as a stream on a channel.
//
// The package also includes a simple semaphore implementation
// that is used internally by the Worker to limit the number of
// concurrent operations.
package hashgrab
