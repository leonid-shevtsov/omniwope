# Discord

Publishes into a [Discord](https://discord.com) channel.

- Create a Discord server (or use an existing one).
- Create a [Discord application](https://discord.com/developers/applications) in the Discord Developer Portal.
  - You will need your own application because you need access to the bot credentials. This is unlike public bots that run on their own servers.
- Create a bot for your application:
  - Go to the "Bot" section in your application
  - Click "Add Bot"
  - Copy the bot token
- Invite the bot to your server:
  - Go to the "OAuth2" > "URL Generator" section
  - Select the "bot" scope
  - Select the following bot permissions:
    - View Channels (required to access channel info)
    - Send Messages
    - Attach Files
    - Embed Links
    - Read Message History
  - Copy the generated URL and open it in your browser to invite the bot
  - **Important**: Make sure the bot has access to the specific channel you want to post to. If the channel has custom permissions, the bot role needs to be allowed in that channel's permission settings.
- Get the channel URL where you want to post:
  - Right-click on the channel in Discord
  - Select "Copy Link" or "Copy Channel Link"
  - This copies a URL like `https://discord.com/channels/123456789/987654321`
- Set `OMNIWOPE_DISCORD_BOT_TOKEN` environment variable to the bot token
- In `omniwope.yml`, configure Discord:
  ```yml
  discord:
    channel: https://discord.com/channels/123456789012345678/987654321098765432
  ```
- Now you are ready to post.

## Configuration

### Channel

Specify the Discord channel using a channel URL:

```yml
discord:
  channel: https://discord.com/channels/123456789012345678/987654321098765432
```

To get the channel URL:

- Right-click on the channel in Discord
- Select "Copy Link" or "Copy Channel Link"
- The URL contains both the server (guild) ID and channel ID

### Start Date

You can set a start date to avoid publishing old posts:

```yml
discord:
  channel: https://discord.com/channels/123456789012345678/987654321098765432
  start_date: 2025-01-01
```

## How It Works

- Posts are published as Discord embeds with rich formatting
- Markdown content is converted to Discord's markdown format
- Images and videos are uploaded as attachments
- Messages can be edited if content changes (detected via checksum)
- Reference links between posts are converted to Discord message links

## Limitations

- Discord has a 2000 character limit for message content
- Embed descriptions have a 4096 character limit (content is truncated if longer)
- Embed titles have a 256 character limit (truncated if longer)
- File uploads are limited to 25MB per file
- Only the first resource (image/video) per post is currently supported

## Permissions Required

The bot needs the following permissions in the target channel:

- **View Channels** (0x0000000400) - Required to access channel information
- **Send Messages** (0x0000000800) - To post messages
- **Attach Files** (0x0000008000) - To upload images and videos
- **Embed Links** (0x00004000) - To use rich embeds
- **Read Message History** (0x00010000) - To edit existing messages

**Combined permissions integer:** `120832` (0x1D800)

You can use this integer value when generating the OAuth2 invite URL programmatically, or select the permissions individually in the Discord Developer Portal.

### Troubleshooting Permission Issues

If you get a "forbidden: bot lacks permissions to access channel" error:

1. **Check server-level permissions**: Make sure the bot was invited with all required permissions listed above.

2. **Check channel-specific permissions**:

   - Go to the channel settings (right-click channel → Edit Channel)
   - Go to the "Permissions" tab
   - Find your bot's role in the list
   - Ensure the bot has all required permissions enabled for that specific channel
   - If the channel has custom permissions that override server permissions, make sure the bot role is explicitly allowed (not denied)

3. **Verify bot is in the server**: The bot must be a member of the server (guild) where the channel is located.

4. **Check role hierarchy**: Make sure the bot's role is not below other roles that might be denying access.

5. **Re-invite the bot**: If permissions were changed after the bot was invited, you may need to re-invite it with the updated permissions.
