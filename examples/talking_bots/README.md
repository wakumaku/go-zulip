Talking Bots is a just for fun example where we have 2 zulip bots sending
messages to each other in a silly conversation.

This example is using the Users, Channels, Messages and Realtime Services.

Outline:

1. Initiate client with User A credentials
2. Initiate client with User B credentials
3. Create necessary services for each one (users, channels, messages and realtime)
4. Start a goroutine for each one where:
    1. They will subscribe to a fixed Channel (#talking_bots)
    2. They will register an Event Queue to only receive notifications on the
    subscribed channel.
    3. They will start sending messages to the channel mentioning or not
    the other user.
    4. They will react or not to the other user messages with another message,
    a reaction or sending a picture.

Running the example:

1. Get 2 user credentials
    * You can run `examples/create_users` to generate 2 random users or create 2 bots
2. Export the variables
    ```sh
    export ZULIP_EMAIL_A="johndoe-bot@localhost"
    export ZULIP_API_KEY_A="xxxxxxxxxxxxxxx"
    export ZULIP_EMAIL_B="janedoe-bot@localhost"
    export ZULIP_API_KEY_B="yyyyyyyyyyyyyyy"
    ```
3. `go run .`
4. Join the channel where the 2 users will talk at https://localhost
