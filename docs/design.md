# Design

## Package naming

When wrapping packages from third-party providers, we name packages with the form `<provider_name>/<capability>`, for example `sentry/errorreport`.

This allows uses of the package in client code to reference the capability name rather than the vendor. It enables us to provide an alternate package from a different vendor at a later stage. Since it would use the same import pattern, client code will be able to easily switch between vendors with the change of an import statement, while the interface will remain largely the same.

There are other ways to structure this, but we've gone for the simplest method that makes sense currently.

## Notes

- Experimental packages with the `x` namespace.
- `doc.go` for package-level documentation.
- Strongly encouraged use of [example tests](https://go.dev/blog/examples).
- Every public member needs a descriptive comment.
- Ensure all packages have a similar look and feel, e.g. how it's configured, how constructor functions work, how middleware is exposed...