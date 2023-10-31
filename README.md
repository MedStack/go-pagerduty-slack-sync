<p align="center">
  <img src="./synclogo.png">
</p>



This tool syncs a pagerduty schedule to one slack group.

For example if you have the following users in a pagerduty schedule:

```
schedule id: 1234
user1
user2
user3 <= currently on call
```

Then when you run a sync for the schedule using:

```
docker run -e RUN_INTERVAL_SECONDS=60 -e SLACK_TOKEN=xxx -e PAGERDUTY_TOKEN=xxx -e SCHEDULE=1234 medstack/pagerduty-slack-sync:latest
```

The following slack group would be created:

- `@on-call` => `user3`


Full parameter list:

| Env Name                | Description                                                                                     | Default Value  | Example                 |
|:------------------------|:------------------------------------------------------------------------------------------------|:---------------|:------------------------|
| PAGERDUTY_TOKEN         | Token used to talk to the PagerDuty API                                                         | n/a            | xxxxx                   |
| SLACK_TOKEN             | Token used to talk to Slack API                                                                 | n/a            | xoxp-xxxxxx             |
| SCHEDULE                | A PagerDuty schedule that you want to sync                                                      | n/a            | 1234                    |
| RUN_INTERVAL_SECONDS    | Run a sync every X seconds                                                                      | 60             | 300                     |
| PAGERDUTY_SCHEDULE_FROM | How far into the future to start the evaluation of Pagerduty schedule (Go time duration format) | 0h             | 8760h                   |                  |


## Slack permissions

In order for the app to run you will need to create a bot with the following permissions:
```
usergroups:read
usergroups:write
users:read
users:read.email
```

If you have locked down your slack so only the admins can create groups then you have two options.  You can either create the slack groups up front and the app will use those or you can give the bot user auth and give it admin perssions:
```
admin.usergroups:read
admin.usergroups:write
usergroups:read
usergroups:write
users:read
users:read.email
```
