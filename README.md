# Standalone Client - Server App

The project's goal is to build a stand-alone client-server application.

The idea is to use a web-based front-end (i.e. an Angular application) for the UX and
a local native application for the business logic.

The project is an experiment to replace a fat client (in this case a swing app) 
exploiting web technologies for the user interaction and a local executable for the 
business logic.

The application is a go program that serves two purposes:

- implement the server side business logic 
- serve the web pages of the client side application
 
As a bonus the go executable will do two other things: 
- open the default browser pointing at it's home page: so that the user just needs to
start the executable 
- listen (on a web socket) for a "ping" from the client: so that when the user closes
the brower the server can terminate

## Build
To build and run the app, first get the external dependencies:

```text
> go get github.com/gorilla/websocket
```

Then build the app:

```text
go build main.go
```

Finally run the app:

```text
./main
```

This will start the server and a browser pointing to the home page of the project.

Web resources will be served from the `static` directory.