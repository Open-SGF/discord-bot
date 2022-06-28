import { config } from 'dotenv';
import { formatEventMessage, fetchTodaysEvent } from './meetupHelpers.js';
import { Client, Intents } from 'discord.js';
import { checkEventIsNovel } from "./discordHelper.js";

config();

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

const sendAnnouncement = async (channel) => {
	const event = await fetchTodaysEvent();
  
  if (!event) {
    console.log('No events scheduled today.');
    return
  }
  
  const timeIsAfterTen = (new Date()).getHours() + 1 > 10;
  
  const eventIsNovel = await checkEventIsNovel(event);
  
  if (!eventIsNovel || !timeIsAfterTen) {
    return
  }
  
	channel.send(formatEventMessage(event));
};

client.on('ready', async () => {
  console.log(`Logged in as ${client.user.tag}!`);
  
  const channel = client.channels.cache.get(process.env.ANNOUNCEMENT_CHANNEL_ID);
  
  await sendAnnouncement(channel);
});

client.login(process.env.LOGIN_TOKEN);
