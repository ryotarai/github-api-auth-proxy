require 'octokit'

Octokit.api_endpoint = 'http://localhost:8080'

client = Octokit::Client.new(login: 'user1', password: 'password')
p client.user
p client.user_issues(state: 'closed')
p client.create_issue('ryotarai/dummy', 'Hello', 'world')
