import fetch from "node-fetch";

export async function checkEventIsNovel (event) {
  const messages = await getAllMessages()
  
  return !messages.some(message => {
    // TODO: see about finding a better way to compare than this url attribute.
      return message.embeds.some(embed => embed.url === event.shortUrl)
    })
}

async function getAllMessages () {
  // TODO: see about limiting the messages that get pulled back.
  const res = await fetch(`https://discord.com/api/channels/${process.env.ANNOUNCEMENT_CHANNEL_ID}/messages`, {
    headers: [['Authorization', `Bot ${process.env.LOGIN_TOKEN}`]]
  })
  
  return await res.json()
}
