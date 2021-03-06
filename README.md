![Lattr](readme/lattr_logo.png)

[![codecov](https://codecov.io/gh/RemeJuan/lattr/branch/main/graph/badge.svg?*token=XeKra2LhuM)](https://codecov.io/gh/RemeJuan/lattr)
[![swagger](https://img.shields.io/badge/-swagger-brightgreen)](http://swagger.lattr.app)

# Lattr

Lattr (later) is a small API driven Twitter scheduler, in practice, it would be something like a lite version
of [Buffer](https://buffer.com) or [Hootsuite](http://hootsuite.com). In the current version, posts can be scheduled
using an API/Webhook, however, a Web/Mobile UI will come a little further down the line.

Webhooks are supported to allow automation with services such as [IFTTT](https://ifttt.com), webhook posts will be
scheduled based on a specified time gap in hours with a randomized minute. You could then post every 2 hours, with some
randomization on the exact minute of posting.

You could also provide a timestamp in the payload to specify a post time.

Additionally, you can configure a time range, too, for example, post only between 6 am and 6 pm, and optionally a daily post
quantity, this would then calculate the hour based on allowed hours divided by total posts to post, thereby
automatically scheduling up additional posts for the following day and potentially building up a long queue.

Things like analytics will not be added in, for that you could use services like [Bitly](http://bit.ly) to shorten your
URL's and track if that way, once a web UI is added I will add in integrations with such services to automate the
process.

## How To...

### Requirements

In order to run the app locally you would need GoLang v1.16+.

To host the project there is a Docker file setup ready to be deployed, outside of that you would need a postgress
database.

In either scenario you would need the following environment variables configured as well as Twitter API credentials.

[Environments and Variables](https://github.com/RemeJuan/lattr/wiki/Environment-Variables)

[Tokens and Scopes](https://github.com/RemeJuan/lattr/wiki/Tokens-and-Scopes)

In `api/tables` you will find the SQL scripts needed to be run to setup the database

Ensure your hosting environment is configured for the correct timezone, all times are relative to the environments
configured timezone

### Running locally

CD into the `api` directory and run `make run`

### Deploying

#### Heroku

Make sure you have the heroku CLI too installed and are logged in.

Run `heroku stack:set container --app APP_NAME`

You can then either add heroku as a remote to your project and push the code up, connect the heroku project to your
GitHUb account for automated deployment.

[How do I set the timezone on my dyno?](https://help.heroku.com/JZKJJ4NC/how-do-i-set-the-timezone-on-my-dyno)

See the [WiKi](https://github.com/RemeJuan/lattr/wiki) for more information and documnetation
