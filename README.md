# Open SGF Discord Bot
This bot's primary purpose is to notify users on Discord about upcoming
events.

[Bot invite link](https://discord.com/api/oauth2/authorize?client_id=945072391950831656&permissions=8590085120&scope=bot)

## Contributing

### Dependencies
- [Golang](https://go.dev/)

### Setup
1. Copy `.example.config.json` to `config.json`
2. Create an Application [Discord Developers](https://discord.com/developers/applications)
3. Add the Bot's token to `$.discordBotToken` in `.config.json`
4. Create a Discord server, and invite the bot to the server

    Todo this, replace the `CLIENT_ID` in the following URL. You can find the 
    Client ID in the Applications -> Your Application -> General Information.

    URL: `https://discord.com/api/oauth2/authorize?client_id=CLIENT_ID&permissions=8590085120&scope=bot`

5. Run the bot with `go run main.go`

### API Documenation
- [Discord's API](https://discord.com/developers/docs/intro)
- [Meetup's API](https://www.meetup.com/api/general/)
- [Discordgo](https://github.com/bwmarrin/discordgo) (Golang Discord library used)

### Get in Touch!:
- [Discord](https://discord.gg/jFD8dZP)
- [Meetup](https://www.meetup.com/open-sgf)

## Deployment/Install

### Initial server setup

1. `docker-compose run debian`
2. Copy `opensgf-discord-bot_*.deb` to the server
4. `sudo dpkg -i /path/to/opensgf-discord-bot_*.deb`
5. `sudo vi /etc/opensgf-discord-bot/env`
6. Write `OPENSGF_DISCORD_BOT_TOKEN=<TOKEN>` (insecure, but only root has access)
7. `sudo systemctl start opensgf-discord-bot.service`

### Performing updates

1. `docker-compose run debian`
2. Copy `opensgf-discord-bot_*.deb` to the server
3. `sudo dpkg -i /path/to/opensgf-discord-bot_*.deb`
4. `sudo systemctl restart opensgf-discord-bot.service`