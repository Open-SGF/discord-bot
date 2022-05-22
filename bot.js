import { config } from 'dotenv';
import { fetchNextEventForGroup, formatEventMessage } from "./meetupHelpers";
import { Client, Intents } from 'discord.js';

config();

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

const sendAnnouncement = async (channel) => {
	const event = await fetchNextEventForGroup(34547654);
 
	channel.send(formatEventMessage(event));
  
  // const nextRun = (time) => {
  //   let timer = setTimeout(async () => {
  //     clearTimeout(timer);
  //     timer = null;
  //     await sendAnnouncement(chan);
  //     timer = nextRun(time);
  //   }, time);
  //
  //   return timer;
  // };
};

client.on('ready', () => {
  console.log(`Logged in as ${client.user.tag}!`);
  
  const channel = client.channels.cache.get(process.env.ANNOUNCEMENT_CHANNEL);

  sendAnnouncement(channel);
});

client.login(process.env.LOGIN_TOKEN);
