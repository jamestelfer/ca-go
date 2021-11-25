// Package flags provides access to feature flags. It wraps the LaunchDarkly SDK
// to expose a convenient and consistent way of configuring and using the
// client.
//
// The client can be configured and used as a managed singleton or as an
// instance returned from a constructor function. The managed singleton provides
// a layer of convenience by removing the need for your application to maintain
// a handle on the flags client.
package flags
