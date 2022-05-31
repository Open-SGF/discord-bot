const { config } = require('dotenv');
const { fetchNextEventForGroup, formatEventMessage } = require("./meetupHelpers");
const { Client, Intents } = require('discord.js');

config();

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

const sendAnnouncement = async (channel) => {
	const event = await fetchNextEventForGroup(34547654);
  
  if (!event) {
    console.log('No upcoming events found.');
    return
  }
 
	channel.send(formatEventMessage(event));
};

client.on('ready', () => {
  console.log(`Logged in as ${client.user.tag}!`);
  
  const channel = client.channels.cache.get(process.env.ANNOUNCEMENT_CHANNEL);

  sendAnnouncement(channel);
});

client.login(process.env.LOGIN_TOKEN);
