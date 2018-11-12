Data Source Interfaces
----------------------

This package contains the golang `interface` that exposes the required functions for Machinable. The purpose of this package is to abstract the database solution used to store the Machinable data. Currently, Machinable only supports a MongoDB a driver, but any database solution could be supported with the correct `dsi` driver.

Refer to the [`Datastore` interface](./interfaces/interfaces.go) if you would like to write a new Machinable DSI driver.


### Errors
Datastore errors are returned as the custom [DatastoreError type](./errors/errors.go). This type implements the `Error` function so it can be used as typical golang `errors`. However, `DatastoreError` also exposts a `Code` function which attempts to translate the "type" of error to a HTTP status code. The purpose of this is to reduce the work the handlers have to do to return an appropriate status code to the user if an error occurs.
