# github-api-auth-proxy

```
+--------+     +-----------------------+      +------------+
|        +---->+                       +----->+            |
| Client |     | github-api-auth-proxy |      | GitHub API |
|        +<----+     w/ OPA policy     +<-----+            |
+--------+     +--------+------+-------+      +------------+
```

## Usage

First, write OPA policy:

```
$ cat <<EOC > policy.rego
package github.authz

default allow = false

allow {
    input.username == "user1"
    input.method == "GET"
    input.path == "/user"
}
EOC
```

In the config file, passwords need to be hashed by bcrypt. You can get bcrypt-ed password as follows:

```
$ github-api-auth-proxy -bcrypt
Password:
Bcrypted: $2a$10$tHUUM6cydLY/Sg9.OqmOsehpRdqbmyAcsdwm6t13qMxAlb4eENur2
```

Second, create a config YAML file:

```yaml
originURL: 'https://api.github.com'
opaPolicyFile: 'policy.rego'
accessToken: 'your personal access token'
passwords:
  user1:
  - '$2a$10$tHUUM6cydLY/Sg9.OqmOsehpRdqbmyAcsdwm6t13qMxAlb4eENur2'
```

Then, start github-api-auth-proxy:

```
$ github-api-auth-proxy -config config.yaml
```

```
$ curl -u user1:password http://localhost:8080/user
```