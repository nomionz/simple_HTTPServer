# Project managment HTTP server

I don't think this is some serious project, rather a cheat sheet for me.

A simple HTTP server. I do this to learn Go. The server will be completely written with the standard packages of the Go language.

The server is a simple version of the project manager. Workers have completed tasks.
-   `GET /workers/{name}` should return the number of the total number of completed tasks
-   `POST /workers/{name}` should record a completed task for that name

Everything is done with the TDD method. That's why in packages you can see some files like '*_test.go'

# How to use it
Build it, run it and then use `curl` to test it out.
- `curl -X POST http://localhost:4673/workers/John` - will append one completed task to the John
- Check tasks with `curl http://localhost:4673/workers/John`
- Get JSON info about project with `curl http://localhost:4673/project`

# Pretty JSON
For pretty output of JSON I'm using `json.SetIndent("", "    ")`
```go
func (p *ProjectManagementServer) projectHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "    ")
    enc.Encode(p.store.GetProjectInfo())
}
```
# Go embedding
For routing purposes I'm using `http.ServeMux`, but actually I can embed `http.Handler` so I shouldn't implement `ServeHTTP`
```go
type ProjectManagementServer struct {
	store ProjectManagementStore
	http.Handler
}
```

# Storage
`fs_store.go` is just a simple file storage implementation.

File is saving in project src directory.