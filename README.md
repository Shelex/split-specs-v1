# Receive spec files to run based on previous run

Simple graphql api to grab spec files in order: new specs -> longest specs -> short specs  
Could be used to make concurrent machines that run your tests much equal in duration times

# Use

- [Minimalistic web UI](https://split-specs.appspot.com) along with
  - [Graphql Playground](https://split-specs.appspot.com/playground)
  - [API Endpoint /query](https://split-specs.appspot.com/query)

# Flow

- Register user or login with existing one
- Create new session (it will be attached to existing project or will create new)
- Get nextSpec for your sessionID and machineID, every query will finish previous spec for this session + machine and return next. Final query will return message "session finished" and finish spec and session for specific machineID. in case machineID is not passed it will be "default"

# Try it locally

- clone this repository
- `cd split-specs`
- `make deps` - download dependencies
- `make keys` - generate private and public keys for auth
- `export ENV=dev` - to use in memory storage instead of real db
- `make api` - build binary and execute
- open `http://localhost:8080/playground` for GraphQL playground
- use `http://localhost:8080/query` for Altair/Postman/Insomnia api clients
- use `http://localhost:8080/` for ui interface

# UI

- `make web-deps` - install dependencies from npm (yarn required)
- `make web-dev` - open ui at localhost for development
- `make web-build` - build page from sources (before deployment)
  As UI is served by golang backend from `web/build` folder, it could be accessed by spawning just backend instance, however without features like hot-reload, thus not so convenient for development

# Use with gcloud sdk

- install Python 2.7, [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) and follow the docs
- [Quickstart](https://cloud.google.com/appengine/docs/standard/go/quickstart) - for using go with appengine
- `make dev` - run dev server with file watcher, app.yaml should have runtime go111, but changed to go115 when deploying
- `make deploy` - deploy app to app engine
- `make browse` - open deployed app in local browser

# Client options

- mutation register - create new user and receive jwt token

```graphql
mutation {
  register(input: { email: "admin@example.com", password: "admin" })
}
```

- mutation login - receive jwt token for existing user

```graphql
mutation {
  login(input: { email: "admin@example.com", password: "admin" })
}
```

- mutation addSession(session) - initialize new session for your project (session will be attached to existing project or create a new one) and receive sessionId

```graphql
mutation {
  addSession(
    session: {
      projectName: "test"
      specFiles: [{ filePath: "1" }, { filePath: "2" }, { filePath: "3" }]
    }
  ) {
    sessionId
    projectName
  }
}
```

- query nextSpec(sessionID, machineID?) - receive next spec file to run for specific session and for specific machineID. In case only one machine is used - no need to pass it

```graphql
query {
  nextSpec(sessionId: "3e1295e4-b044-4a7a-82a7-b0e71afe70e7")
}
```

- query project(name): get your project and all sessions info

```graphql
query {
  project(name: "test") {
    projectName
    sessions {
      id
      start
      end
      backlog {
        file
        estimatedDuration
        start
        end
      }
    }
  }
}
```

- query session(sessionId): get session by id

```graphql
query {
  session(sessionId: "3e1295e4-b044-4a7a-82a7-b0e71afe70e7") {
    id
    start
    end
    backlog {
      file
      start
      end
      estimatedDuration
      assignedTo
    }
  }
}
```

- query projects: get list of project names available for current user

```graphql
query {
  projects
}
```

- mutation shareProject: make your project available for other existing user

```graphql
mutation {
  inviteUser(email: "admin2", projectName: "test")
}
```

- mutation changePassword: change password for signed in user

```graphql
mutation {
  changePassword(input: { password: "admin", newPassword: "ababagalamaga" })
}
```

- mutation deleteSession: delete existing session by id for current user

```graphql
mutation {
  deleteSession(sessionId: "vcV8iLiN_Z5rEsMlF8ur1")
}
```

- mutation deleteProject: delete existing project by name for current user

```graphql
mutation {
  deleteProject(projectName: "test")
}
```
