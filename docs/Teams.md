# Teams Approach

**Contents**

- [Description](#description)
- [Vocabulary](#vocabulary)
- [Approach](#approach)
  - [Teams](#teams)
  - [API Resources](#api-resources)
  - [Key/Value Store](#key-value-store)
  - [Webhooks](#webhooks)
  - [Logs](#logs)

## Description

Project API Resources and Key/Value Store (and their objects) currently have Authorization settings only relevant to individual project users. This is useful in some scenarios, but for _most_ real-world applications, the concept of Groups/Organizations/Teams exists where multiple users share and manage access to a set of resources. This document attempts to outline the approach, design, and changes necessary to support Teams of Project Users, as well as their authorization policies to a set of objects.

## Vocabulary

To ensure any interested can agree on and understand this document, this section outlines the vocabulary of words and phrases used to outline the approach described above. Some of the defined vocabulary will be repetitive compared to the [Machinable User Documentation](https://www.machinable.io/documentation/).

| Word              | Definition                                                                                                                                                                 |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Resources**     | Any Machinable API endpoints (API Resources, Key/Value store, Teams). More specifically, any API Endpoints that require Project level authorization checks                 |
| **Authorization** | Specifying and verifying Project User access rights/priveleges to a Resource                                                                                               |
| **Users**         | For the purpose of this document, Users refers to **Project Users**. **Project Users** are clients that authenticate and make requests to a Machinable Project's Resources |
| **Team**          | A team is a logical grouping of Users that influences Authorization to a subset of Resources and data stored for Resources                                                 |

> _NOTE: To clarify, this document only refers to the implementation of Project Teams (teams of users of a Machinable project), not Application Teams(teams of users within Machinable itself). Machinable Teams will be a separate development effort that allows Machinable Users to manage teams and access to projects at an adminstrative level_

---

## Approach

### Teams

Project Teams will be a new entity within a Machinable Project that provides the necessary API Endpoints to manage a Team and its Users. Any project user will be able to create a team (based on project settings). The initial creator of a Team will be its first `administrator`. Team Administrators have full access to invite and remove users from a team.

#### Schema

**Team**

| Field       | Type       | Description                                                 |
| ----------- | ---------- | ----------------------------------------------------------- |
| ID          | `uuid`     | The unique identifier of a team                             |
| Name        | `string`   | The human-readable name of a team                           |
| Slug        | `string`   | Unique URL slug of the team (generated from name)           |
| Description | `string`   | Longer description of the team, it's purpose, etc.          |
| Created     | `datetime` | The timestamp of Team creation                              |
| Image       | `string`   | A URL of an avatar for this team                            |
| Creator     | `uuid`     | The unique identifier of the User that created the team(FK) |

**Team Member**

| Field      | Type       | Description                                                                                                           |
| ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------- |
| Team ID    | `uuid`     | The team's ID (fk)                                                                                                    |
| User ID    | `uuid`     | The user's ID (fk)                                                                                                    |
| User Email | `string`   | The email of the team member. This is needed for initial team invites.                                                |
| Joined     | `datetime` | The timestamp of when the user joined the team                                                                        |
| Role       | `enum`     | The role of this team member: `admin`, `member`, `observer`, more? This influences a users access and priveleges      |
| Status     | `enum`     | The status of this users membership to the team: `active`, `invited`, `inactive`. This influences what a user can do. |

#### API

Project Team API Endpoints and their purpose.

**Hostname and Path Prefix**: `https://{project-slug}.machinable.io`

| Name         | Path                                   | Verb     | Description                                                      | Authz Notes                                                                 |
| ------------ | -------------------------------------- | -------- | ---------------------------------------------------------------- | --------------------------------------------------------------------------- |
| CreateTeam   | `/teams`                               | `POST`   | Creates a new team                                               | Creates a team with the creator as the first and single (admin) team member |
| GetTeams     | `/teams`                               | `GET`    | Retrieves all teams for which the requesting user is a member    | Requestor's teams                                                           |
| GetTeam      | `/teams/:teamIdOrSlug`                 | `GET`    | Retrieve the details of a team and the list of its team members  | Retrieve team by ID and Requestor User ID as a team member (view)           |
| UpdateTeam   | `/teams/:teamIdOrSlug`                 | `PUT`    | Updates a team's information; description and/or image           | Verify requestor is admin of team                                           |
| DeleteTeam   | `/teams/:teamIdOrSlug`                 | `DELETE` | Permanently delete a team. Team will need to have 0 team members | Verify requestor is admin of team                                           |
| InviteMember | `/teams/:teamIdOrSlug/invite`          | `POST`   | Invites a new user, by email, to join the Team                   | Verify requestor is admin of team                                           |
| ListMembers  | `/teams/:teamIdOrSlug/members`         | `GET`    | Lists all users of a team                                        | Verify requestor is a member of the team                                    |
| DeleteMember | `/teams/:teamIdOrSlug/members/:userId` | `DELETE` | Permanently removes a member of the team                         | Verify requestor is admin of team                                           |

#### Data Flows

- `InviteMember`
  1. Create `Team Member` record with invite email and `status` set to `invited`
  2. Send email to the invited user with time sensitive invite link (we cannot assume they are already registered as a user of the project)
  3. User signin/signup, follows invite URL (redirect), is added to team.
  4. Set `Team Member` record with `user_id`, `joined`, and `status` => `active`
- `DeleteMember`
  1. Remove `Team Member` record
  2. Send email to the user to inform them they were removed

### API Resources

In addition to the team API Paths documented above, Team objects stored for an API Resource will be managed at new API Endpoints.

**Hostname and Path Prefix**: `https://{project-slug}.machinable.io`

| Name         | Path                                               | Verb     | Description                                 | Authz Notes                                                                                                              |
| ------------ | -------------------------------------------------- | -------- | ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| ListObjects  | `/teams/:teamIdOrSlug/api/:resourceSlug`           | `GET`    | Retrieve the list of objects for a resource | Verify requestor is a member of the team                                                                                 |
| GetObject    | `/teams/:teamIdOrSlug/api/:resourceSlug/:objectId` | `GET`    | Retrieve a single object                    | Verify requestor is a member of the team                                                                                 |
| CreateObject | `/teams/:teamIdOrSlug/api`                         | `POST`   | Create a single object                      | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `create` privelege) |
| PutObject    | `/teams/:teamIdOrSlug/api/:resourceSlug/:objectId` | `PUT`    | Update a single object                      | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `update` privelege) |
| DeleteObject | `/teams/:teamIdOrSlug/api/:resourceSlug/:objectId` | `DELETE` | Delete a single object                      | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `delete` privelege) |

**Optionally**, a request header could be provided to the _existing_ resource API Endpoints that identifies a user's intention of managing a _team's_ resource objects. In this case, the proposed header would be:

| Header          | Description                                                                    |
| --------------- | ------------------------------------------------------------------------------ |
| `X-Mchn-TeamID` | Request header used to identify a Team the user would like to manage Resources |

### Key/Value Store

Currently, authorization of Keys is limited to authentication options only. Implementing true authorization (creator being the owner, etc) should be a separate effort.

Team resources are globally accessible to the Team Members, i.e. they are _at least_ readable by any team member. This will be the case for Keys as well; Root keys will be bucketed to a team by `teamId`, with access being dictated by the Team Member Role/Access settings. As with API Resources, Key/Value Store will have new API Endpoints and, **optionally**, a request header to be used with the current endpoint.

**Hostname and Path Prefix**: `https://{project-slug}.machinable.io`

| Name        | Path                                    | Verb     | Description                         | Authz Notes                                                                                                              |
| ----------- | --------------------------------------- | -------- | ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| GetKey      | `/teams/:teamIdOrSlug/json/:rootKey`    | `GET`    | Retrieve a single key               | Verify requestor is a member of the team                                                                                 |
| GetKeyValue | `/teams/:teamIdOrSlug/json/:rootKey/**` | `GET`    | Retrieve the value at the JSON path | Verify requestor is a member of the team                                                                                 |
| CreateKey   | `/teams/:teamIdOrSlug/json`             | `POST`   | Create a key                        | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `create` privelege) |
| PutValue    | `/teams/:teamIdOrSlug/json/:rootKey/**` | `PUT`    | Update value at the JSON path       | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `update` privelege) |
| DeleteValue | `/teams/:teamIdOrSlug/json/:rootKey/**` | `DELETE` | Delete value at the JSON path       | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `delete` privelege) |
| DeleteKey   | `/teams/:teamIdOrSlug/json/:rootKey`    | `DELETE` | Delete root key                     | Verify requestor is a member of the team and has the appropriate team role (`admin` or `member` with `delete` privelege) |

**Optionally**, a request header could be provided to the _existing_ resource API Endpoints that identifies a user's intention of managing a _team's_ resource objects. In this case, the proposed header would be:

| Header          | Description                                                                     |
| --------------- | ------------------------------------------------------------------------------- |
| `X-Mchn-TeamID` | Request header used to identify a Team the user would like to manage Key/Values |

### Webhooks

Include team ID in webhooks.

Include webhooks for team events:

1. Team member joins a team
2. Team member is removed from a team

### Logs

Ensure current logging middleware is used with new team HTTP endpoints.

### Settings

- Enable/disable teams at a project level
- Enable/disable teams at a resource/key level
