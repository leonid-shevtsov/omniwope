# Mastodon

Publishes into a [Mastodon](https://joinmastodon.org) account.

- Create an account.
- Obtain authorization token - you can use [GetAuth for Mastodon](https://getauth.thms.uk/?scopes=read+write:statuses&client_name=OmniWOPE).
- Set `OMNIWOPE_MASTODON_ACCESS_TOKEN` environment variable to the access toekn
- In `omniwope.yml`, set the instance URL and other configuration options:
  ```yml
  mastodon:
    instance_url: https://mastodon.social
    visibility: public
    language: en
  ```
- Now you are ready to post.

## A note about the Fediverse

_Technically,_ you are not publishing "into Mastodon". You are publishing into the Fediverse! But this output is called Mastodon, because it's using the Mastodon client API - which has become a de-facto standard for microblog-flavored Fediverse servers like GoToSocial or Pleroma.

## A note about ActivityPub

Technically, your blog can itself be an ActivityPub enabled server, and then you don't need a separate account to republish there. However, this solution is much more complicated than just OmniWOPing to a Mastodon account of your choice.
