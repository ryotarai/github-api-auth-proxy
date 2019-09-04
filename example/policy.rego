package httpapi.authz

# HTTP API request
import input

default allow = false

allow {
    input.header["Authorization"][_] == "token hello"
    input.method == "GET"
    input.path == "/user"
}

allow {
    input.header["Authorization"][_] == "token hello"
    input.method == "GET"
    input.path == "/user/issues"
    input.query["state"][_] == "closed"
}

