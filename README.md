# Split your specs based on previous run durations
Simple graphql api to grab spec files in order: new specs -> longest specs -> short specs  
Could be used to make parallel machines that run your tests much equal in duration times  
Deployed with Google Cloud AppEngine

# Use
 - [Graphql Playground](https://split-specs.appspot.com/)
 - [API Endpoint /query](https://split-specs.appspot.com/query)


# Install
 - clone this repository  
 - `cd split-specs`  
 - `make deps` - download dependencies
 - `make api` - build binary and execute  
 - open `http://localhost:8080/` for GraphQL playground or use Altair/Postman/Insomnia  

# Develop with gcloud sdk
 - `make dev` - run dev server with file watcher, app.yaml should have runtime go111
 - `make deploy` - deploy app to app engine
 - `make browse` - open deployed app in local browser

# Client options
  - mutation register - create new user and receive jwt token
```graphql
mutation{
  register(input: {
    username: "admin",
    password: "admin"
  })
}
 ```

   - mutation login - receive jwt token for existing user
```graphql
mutation{
  login(input: {
    username: "admin",
    password: "admin"
  })
}
 ```

 - mutation addSession(session) - initialize new session for your project (session will be attached to existing project or create a new one) and receive sessionId
```graphql
 mutation {
  addSession (session:{
    projectName: "test",
    specFiles: [
      {
        filePath: "1"
      },
      {
        filePath: "2"
      },
      {
        filePath: "3"
      }
    ]
  }) {
    sessionId
    projectName
  }
}
 ```

 - query nextSpec(sessionID, machineID?) - receive next spec file to run for specific session and for specific machineID. In case only one machine is used - no need to pass it
```graphql
query {
  nextSpec (sessionId: "3e1295e4-b044-4a7a-82a7-b0e71afe70e7")
}
 ```

  - query project(name): get your project and all sessions info
```graphql
query{
  project(name: "test"){
    projectName
    latestSession
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

  - query projects: get list of project names available for current user
```graphql
query{
  projects
}
```

  - mutation shareProject: make your project available for other existing user
```graphql
mutation{
  inviteUser(username: "admin2", projectName: "test")
}
```

  - mutation changePassword: change password for signed in user
```graphql
mutation{
  changePassword(input: {
    password: "admin",
    newPassword: "ababagalamaga"
  })
} 

```

 # TODO
  - :x: implement persistance layer with firestore in datastore mode
  - :x: implement detection logic for changes in test list inside spec file