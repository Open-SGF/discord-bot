import fetch from "node-fetch";

export async function checkEventIsNovel (event) {
  const messages = await getAllMessages()
  
  return !messages.some(message => {
      return message.embeds.some(embed => embed.url === event.shortUrl)
    })
}

async function getAllMessages () {
  const res = await fetch(`https://discord.com/api/channels/${process.env.ANNOUNCEMENT_CHANNEL_ID}/messages`, {
    headers: [['Authorization', `Bot ${process.env.LOGIN_TOKEN}`]]
  })
  
  return await res.json()
}
