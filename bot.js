import { config } from 'dotenv';
import { fetchNextEventForGroup, formatEventMessage } from "./meetupHelpers";
import { Client, Intents } from 'discord.js';

config();
const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

const sendAnnouncement = async (channel) => {
	const response = await fetchNextEventForGroup(process.env.MEETUP_GROUP_ID);
	channel.send(formatEventMessage(response));
};

client.on('ready', () => {
  console.log(`Logged in as ${client.user.tag}!`);
  
  const channel = client.channels.cache.get(process.env.ANNOUNCEMENT_CHANNEL);

  sendAnnouncement(channel);
});

client.login(process.env.LOGIN_TOKEN);


