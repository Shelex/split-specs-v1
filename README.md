# Split your specs based on previous run durations
Simple graphql api to grab spec files in order: new specs -> longest specs -> short specs  
Could be used to make parallel machines that run your tests much equal in duration times  
Deployed with Google Cloud AppEngine

# Use
 - [Graphql Playground](https://test-splitter.appspot.com/)
 - [API Endpoint /query](https://test-splitter.appspot.com/query)


# Install
 - clone this repository  
 - `cd split-test`  
 - `make deps` - download dependencies
 - `make api` - build binary and execute  
 - open `http://localhost:8080/` for GraphQL playground or use Altair/Postman/Insomnia  

# Develop with gcloud sdk
 - `make dev` - run dev server with file watcher, app.yaml should have runtime go111
 - `make deploy` - deploy app to app engine
 - `make browse` - open deployed app in local browser

# Client options

 - mutation addSession(session) - initialize new session for your project and receive sessionId
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

 - query nextSpec(sessionID, machineID?) - receive next spec file to run for specific session and possibly for specific machineID. In case only one machine is used - no need to pass it
```graphql
query {
  nextSpec (sessionId: "3e1295e4-b044-4a7a-82a7-b0e71afe70e7")
}
 ```

  - query project(name): get your project and sessions info
```
query {
  project(name:"test") {
    projectName
    latestSession
    sessions{
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


 # TODO
  - :x: implement persistance layer with firestore in datastore mode
  - :x: implement detection logic for changes in test list inside spec file