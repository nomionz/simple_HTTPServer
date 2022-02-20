# Project managment HTTP server

A simple HTTP server. I do this to learn Go. The server will be completely written with the standard packages of the Go language.

The server is a simple version of the project manager. Workers have completed tasks.
-   `GET /workers/{name}` should return the number of the total number of completed tasks
-   `POST /workers/{name}` should record a completed task for that name

# How to use it
Build it, run it and then use `curl` to test it out.
- `curl -X POST http://localhost:4637/workers/John` - will append one completed task to the John
- Check tasks with `curl http://localhost:4637/workers/John` 

# Feature updates
- GET JSON
- Store the variables to the file
