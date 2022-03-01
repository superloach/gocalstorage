# gocalstorage
Go bindings for the JavaScript Storage APIs, using WASM.

## tests
Tests use https://github.com/agnivade/wasmbrowsertest.

StorageEvent tests are currently fragile, using an iframe hack to simulate
multiple pages on the same domain.
