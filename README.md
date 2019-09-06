# github-api-auth-proxy

```
+--------+     +-----------------------+      +------------+
|        +---->+                       +----->+            |
| Client |     | github-api-auth-proxy |      | GitHub API |
|        +<----+                       +<-----+            |
+--------+     +--------+------+-------+      +------------+
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
$ opa run -s policy.rego
```

Second, create a config YAML file:

```yaml
originURL: 'https://api.github.com'
opaServerURL: 'http://localhost:8181'
accessToken: 'your personal access token'
passwords:
  user1:
  - '$2a$12$j0mFVW4Ccx3AV8I1URRq.uB.7cziHxnLdIrSjrpRcwJKfDLEw410W'
```

In the config file, passwords need to be hashed by bcrypt. You can get bcrypt-ed password as follows:

```
$ github-api-auth-proxy -bcrypt
Password:
Bcrypted: $2a$10$L1iwhyLeVWmy5Rj1zu1lzujT4S4KhKUEFaiLUzg7NzsXcMvNAT7t2
```

Then, start github-api-auth-proxy:

```
$ github-api-auth-proxy -config config.yaml
```
