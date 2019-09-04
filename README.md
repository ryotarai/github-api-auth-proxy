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

```
github-api-authz-proxy -origin-url https://api.github.com -access-token 'Your Personal Access Token' -opa-server-url 'OPA server'
```