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

# Feature updates
- GET JSON
- Store the variables to the file
