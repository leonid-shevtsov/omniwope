# Telegram

Publishes into a [Telegram](https://telegram.org) channel.

- Create a channel.
- Set up a [bot](https://core.telegram.org/api).
- Add bot to channel administrators.
- Set `OMNIWOPE_TG_CREDENTIALS` environment variable to the bot credentials
- In `omniwope.yml`, set the channel name:
  ```yml
  tg:
    channel: mychannel_name
  ```
- Now you are ready to post.
