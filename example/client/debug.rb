require 'octokit'

Octokit.api_endpoint = 'http://localhost:8080'

client = Octokit::Client.new(access_token: 'hello')
p client.user
p client.user_issues(state: 'closed')
p client.create_issue('ryotarai/dummy', 'Hello', 'world')
