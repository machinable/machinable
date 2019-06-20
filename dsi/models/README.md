The models in this package are still using fields specific to a MongoDB backend (particularly ObjectID). In order to _really_ support "any" datastore, these models need to be completely agnostic (only use golang types).

The following files contain the `ObjectID` type in their models:

* `access.go`
* `projects.go`
* `users.go`
* `sessions.go`
