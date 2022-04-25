require('dotenv').config();
const { Client, Intents } = require('discord.js');
const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

const fetchMeetupData = async () => {
	// fetch stuff here
};

const sendAnnouncement = async (channel) => {
	await fetchMeetupData();
	channel.send('data fetched');
};

client.on('ready', () => {
  console.log(`Logged in as ${client.user.tag}!`);
  const chan = client.channels.cache.get(process.env.ANNOUNCEMENT_CHANNEL);

  // Poll for meetup data every day, if there's an event that day
  // then go ahead and make the announcement

  sendAnnouncement(chan);

  const nextRun = (time) => {
    let timer = setTimeout(async () => {
       clearTimeout(timer);
       timer = null;
       await sendAnnouncement(chan);
       timer = nextRun(time);
    }, time);

    return timer;
  };

  nextRun(3 * 1000);
});

client.login(process.env.LOGIN_TOKEN);
