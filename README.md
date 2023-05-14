A drone plugin to send build status messages to your selected ntfy server.

# Usage

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: ntfy notification
  image: monster010/drone-ntfy
  pull: if-not-exists
  when:
    status: [success, failure]
  settings:
    url: https://ntfy.example.org
    topic: events
    priority: high
    tags:
      - pipeline-status
      - dev
    token:
      from_secret: ntfy_token
```

## Properties

`url` *string* [optional] \
Ntfy server.
> *Default: https://ntfy.sh*

`topic` *string* [**REQUIRED**] \
Topic to publish message.
> *Default: none*

`priority` *string* [optional] \
Priority of the notification. Values can be [min, low, default, high, max].
> *Default: default*

`tags` *string* [optional] \
Custom tags to include.
> *Default: none*

`username` *string* [optional] \
Username with publish permissions.
> *Default: none*

`password` *string* [optional] \
[***SECRET RECOMMENDED***] \
Password for username.

> *Default: none*

`token` *string* [optional] \
[***SECRET RECOMMENDED***] \
Token to use, instead username and password.

> *Default: none*

