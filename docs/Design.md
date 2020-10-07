# Machinable

**Contents**

* [Users](#users)
* [Teams](#teams)
* [Projects](#projects)
  * [API Resources](#api-resources)
  * [Collections](#collections)
  * [Access](#access)
  * [Security](#security)
  * [Settings](#settings)

## Users

Machinable users are admins of the projects created. A user can be part of many teams.

Standard user registration/login with JWT and refresh tokens. UI uses the user accessible API to manage projects and teams.

## Teams

Machinable teams are groups of users that manage team projects. A team can have many users.

## Projects

A machinable project is a project that is created and managed by a machinable user or team. A project will have a unique hostname, `{project-slug}.machinable.io`, that will provide access to the projects API Resource definition, collection, and auth endpoints.

The API Resources and collections are HTTP interfaces to a NoSQL database.

There are many purposes of a project:

* Prototype HTTP API
  * Defined and validated models
* Dummy HTTP API
  * Defined and validated models
* Unstructered JSON Storage in a collection
  * This has many uses, including the two above, but without defining any models/resources on the project/server
* Complete CRUD and User access for a production application (reach goal)
  * With a defined API and/or unstructured JSON storage, User management with authz settings, machinable projects could be used to standup a production service quickly.
  * More thought and research needed for this
    * Webhooks on CRUD operations
    * User confirmation (email)
    * User password reset (email)
    * Websockets for updates to data
  

### API Resources

**Overview**

Resources will provide a user defined HTTP API with defined endpoints that correspond to models, defined using a subset of [JSON Schema](https://json-schema.org/). By defining the models with JSON Schema, the objects that are created using the API can be validated. User submitted data will be validated against the defined schema using the [Go OpenAPI libraries](https://github.com/go-openapi).

For example, a defined resource may look like:

```json
{
  "title": "Dogs",
  "path_name": "dogs",
  "properties": {
    "age": {
      "type": "integer", 
      "description": "The age of the dog."
    },
    "birthdate": {
      "type": "string", 
      "format": "date-time", 
      "description": "The birth date of the dog."
    },
    "name": {
      "type": "string", 
      "description": "The name of the dog."
    },
    "breed": {
      "type": "string", 
      "description": "The breed of the dog."
    },
    "commands": {
      "type": "array", 
      "items": {
        "type": "string"
      }, 
      "description": "A list of commands that the dog knows."
    }
  }
}
```

If our project for this resource was called `pets`, our hostname to access this resource would look like `pets.mchbl.com`. After defining and creating the above resource, we could perform CRUD operations on it with:

`GET https://pets.mchbl.com/api/dogs`

`POST https://pets.mchbl.com/api/dogs`

`DELETE https://pets.mchbl.com/api/dogs`

`GET https://pets.mchbl.com/api/dogs/{id}`

`PUT https://pets.mchbl.com/api/dogs/{id}`

`DELETE https://pets.mchbl.com/api/dogs/{id}`

The collection of `dogs` will be returned as the payload:

```json
{
  "items": [
    {
      ...
    },
    ...
  ]
}
```

**Types**

Supported property types:

* `string`
* `integer`
* `number`
* `array`
  * `array` types will only support `string`, `integer`, and `number`
* `object`
  * `object` types should not support nested `objects`, to reduce complexity (or _some_ defined amount of nesting)

**Formats**

Supported formats:

* `date-time`
  * RFC3339 date time `string`
  
**UPDATE**: By using the go-openapi validation libraries, Machinable will support all formats of properties.

**Status Codes**

Possible return status codes for a defined API Resource

* `200 OK`
* `201 Created`
* `400 Bad Request`
* `401 Unauthorized`
* `404 Not Found`
* `500 Internal Server Error`

More TBD.

**Querying**

Filtering on a list by property, for the `dogs` example resource, could look like:

```
cURL -s https://pets.mchbl.com/api/dogs?breed=labrador&age=10
{
  "items": [
    ...
  ]
}
```

i.e. any existing, primitive, fields should be queryable with the `=` operator.

Future:

* Pagination and other comparison operators (`>`, `<`, etc.)
* Filter on objects/arrays

**Access**

Set access policy per resource (or global to the project?).

* Open Policy
  * Anyone who has the URL to the resource can access it
* User/Token Policy
  * Access is dictated by the project Access rules (users/tokens)

_**NOTE**: A project user/token is defined by the project and can only access the defined resources, collections, and project auth endpoints. A Machinable user/team is an admin to the project and can manage the project from the UI._

_**NOTE**: The machinable user/team account will be able to read/write to all resources for manageability. Perhaps this can be configured per user of the team._

**MongoDB Design**

Mongodb collection naming:

* `prj.{project-slug}.definitions` - API resource definitions
* `prj.{project-slug}.api.{path_name}` - API objects per resource

**Other Thoughts**

* Relations
  * Supporting native relations by setting a `key` name on the resource definition, then other resources could use that `key` as a `type`.
  * The `property` with the `key` `type` would store the `_id` of the keyed resource, we can then translate relations
  * Users could also just document their own relations with the current design...
* Validation keywords
  * Support all other type validation keywords - 
Future: Support different type criteria - https://json-schema.org/understanding-json-schema/reference/
* Import/export swagger
  * Once all validation keywords are supported, we could import/export swagger files
* Data visualization
  * Create custom dashboards/visualizations of your data

### Collections

Collections provide a direct, but limited, HTTP interface to a set of defined MongoDB collections. Once a collection is created by a Machinable User/Team, any JSON document can be stored there, i.e. JSON structure is not enforced. Items of a collection are returned in the same format as an API Resource:

```json
{
  "items": []
}
```

Depth of the collection items must be limited. Initially, a depth of `8` seems reasonable, anything beyond that hints at a flawed JSON schema design.

Items of a collection can be accessed and edited in the same way as an API Resource:

With a collection called `settings`:

`GET https://pets.mchbl.com/collections/settings`

`POST https://pets.mchbl.com/collections/settings`

`DELETE https://pets.mchbl.com/collections/settings`

`GET https://pets.mchbl.com/collections/settings/{id}`

`PUT https://pets.mchbl.com/collections/settings/{id}`

`DELETE https://pets.mchbl.com/collections/settings/{id}`

**Querying**

Filtering is similar to API resources. Since the JSON documents are unstructured, be aware of what your documents _could_ look like, as this will impact filtering.

**Status Codes**

Status codes will be the same as API resources.

**Access**

Set access policy per collection (or global to the project?).

* Open Policy
  * Anyone who has the URL to the collection can access it
* User/Token Policy
  * Access is dictated by the project Access rules (users/tokens)
  
**MongoDB Design**

Mongodb collection naming:

* `{project-slug}.collections`
* `{project-slug}.collections.{collection}`

**Other Thoughts**

* Should collections just be a JSON tree instead of a list of items? Perhaps an option for both?
  * User defined `root` type: `object` or `array` of objects?
* Data visualization
  * Create custom dashboards/visualizations of your data

### Access

**Users**

_NOTE: Project users can only access the resource, collection, authn endpoints. They cannot use the management APIs._

User registration/authentication with username/email + password. Admins can configure what resources/collections each user can read/write.

Users can authenticate to the project's `/project/sessions` endpoint, which will return a JWT. JWTs are valid for 5 minutes and can be refreshed using the refresh token. The refresh token is valid for `x` days and represents a user's session, which is viewable in the `Security` section. The JWT expiration (access and refresh) should be configurable per project.

**Tokens**

API Tokens are generated strings that can be used to directly interact with the resource and collection endpoints.

### Security

Security allows the Machinable users/team to view resource/collection activity logs as well as active user sessions. Active user sessions can be revoked.

### Settings

Settings provides an interface to configure the project team (enable/disable required authentication) as well as billing.

---
&copy; Nick Sjostrom
