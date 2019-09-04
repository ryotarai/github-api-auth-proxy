# github-api-authz-proxy

```
+--------+     +------------------------+      +------------+
|        +---->+                        +----->+            |
| Client |     | github-api-authz-proxy |      | GitHub API |
|        +<----+                        +<-----+            |
+--------+     +--------+------+--------+      +------------+
                        |      ^
                        |      |
                        |      |
                        v      |
                     +--+------+--+
                     |            |
                     | OPA server |
                     |            |
                     +------------+
```

## Usage

First, start OPA server ([example policy](example/policy.rego)):

```
opa run -s policy.rego
```

Then, start github-api-authz-proxy:

```
github-api-authz-proxy -origin-url https://api.github.com -access-token 'Your Personal Access Token' -opa-server-url 'OPA server'
```
