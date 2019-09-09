package github.authz

default allow = false

allow {
    input.username == "user1"
    input.method == "GET"
    input.path == "/user"
}

allow {
    input.username == "user1"
    input.method == "GET"
    input.path == "/user/issues"
    input.query["state"][_] == "closed"
}

