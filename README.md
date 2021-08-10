![Lattr](readme/lattr_logo.png)

# Lattr

Lattr (later) is a small API driven twitter scheduler, in practice it would be something like a lite version
of [Buffer](https://buffer.com) or [Hootsuite](http://hootsuite.com). In the current version posts can be scheduled
using an API/Webhook, however a Web/Mobile UI will come a little further down the line.

Webhooks are supported to allow automations with services such as [IFTTT](https://ifttt.com), webhook posts will be
scheduled based on a specified time gap in hours with a randomized minute. You could then post every 2 hours, with some
randomization on the exact minute of posting.

You could also provide a time stamp in the payload to specify a post time.

Additionally, you can configure a time range, to for example, post only between 6am and 6pm, and optionally a daily post
quantity, this would then calculate the hour based on allowed hours divided by total posts to post, thereby
automatically scheduling up additional posts for the following day and potentially building up a longer queue.

Things like analytics will not be added in, for that you could use services like [Bitly](http://bit.ly) to shorten your
URL's and track if that way, once a web ui is added I will add in integrations with such services to automate the
process.

# How To...

Coming soon, this is still in early development.
