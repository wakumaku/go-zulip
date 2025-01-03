Create users for testing via API

```sh
export ZULIP_EMAIL="admin@zulip.org"
export ZULIP_API_KEY="123456"
export ZULIP_SITE="https://localhost"

go run . -email "talkbot1@zulip.org" -password "123456" -name "Talker 1"
```
