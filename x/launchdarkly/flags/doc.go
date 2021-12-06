// Package flags provides access to feature flags. It wraps the LaunchDarkly SDK
// to expose a convenient and consistent way of configuring and using the
// client.
//
// The client can be configured and used as a managed singleton or as an
// instance returned from a constructor function. The managed singleton provides
// a layer of convenience by removing the need for your application to maintain
// a handle on the flags client.
//
// To configure the client as a singleton:
//   err := flags.Configure(
//	   flags.WithSDKKey("foobar"),
//   )
//   if err != nil {
//     // handle invalid configuration
//   }
//
//   err = flags.Connect()
//   if err != nil {
//	   // handle errors connecting to LaunchDarkly
//	 }
//
// To configure the client as a instance that you manage:
//   client, err := flags.NewClient(
//	   flags.WithSDKKey("foobar"),
//   )
//   if err != nil {
//	   // handle invalid configuration
//	 }
//
// The client can be configured to use the Relay Proxy in either Daemon or Proxy
// modes. See WithDaemonMode and WithProxyMode for more information.
//
// Querying for flags is done on the client instance. You can get instance from the
// managed singleton with GetDefaultClient():
//   client, err := flags.GetDefaultClient()
//	 if err != nil {
//	   // client not configured or connected
//	 }
//
// Query a flag by supplying the request context, the flag name, and fallback value
// to be used if an error occurs. The fallback value will always be reflected as
// the value of the flag if err is not nil. The SDK will attempt to extract request
// fields and the authenticated user from the context.
//   val, err := client.QueryBool(ctx, "my-flag", false)
//
// You can also supply your own user object instead of a context:
//   user := flags.NewUser(
//			   "user-id",
//			   flags.WithCustomerAccountID("account-id"),
//	 )
//   // user := flags.AnonymousUser() // if unauthenticated
//
//   val, err := client.QueryBoolWithUser("my-flag", user, false)
//
// When your application is shutting down, you should call Shutdown() to gracefully
// close connections to LaunchDarkly:
//   client.Shutdown()
package flags
