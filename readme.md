# TodoMVC in Treetop

Fully server side application with all the interactivity of a client MV* application.

## Fully server side?

There is no client side templating code of any sort. All IO is driven through XHR and HTML fragments.

See [Treetop](https://github.com/rur/treetop) for more details.

App template by [TodoMVC App Template](https://github.com/tastejs/todomvc-app-template)

## Build Instructions

Assuming you have a `$GOPATH` setup and `npm` is installed.

    go get github.com/rur/todowithoutmvc
    cd $GOPATH/src/github.com/rur/todowithoutmvc/
    npm install
    go run start.go

After the server is running go to http://localhost:8000 and starting using the app. Activate your network inspector to help see what is going on.

_Note, this app makes use of the POST redirect GET pattern via XHR which is opaque to many network inspectors._
