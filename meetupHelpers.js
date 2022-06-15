import fetch from 'node-fetch';
import { joinArrayHumanReadable } from './utils.js';

export async function fetchTodaysEvent () {
  const nextEvent = await fetchNextEventForGroup(34547654)
  
  const eventDate = nextEvent.dateTime.split('T')[0]
  
  const now = new Date();
  
  const todayArray = [now.getFullYear(), now.getMonth() + 1, now.getDate()]
  
  if (todayArray[1] < 10) {
    todayArray[1] = `0${todayArray[1]}`
  }
  
  if (todayArray[2] < 10) {
    todayArray[2] = `0${todayArray[2]}`
  }
  
  if (eventDate !== todayArray.join('-')) {
    return null
  }
  
  return nextEvent
}

export async function fetchNextEventForGroup (groupId) {
  const variables = { groupId }
  
  const query = `query GetUpcomingEventsForGroup ($groupId: ID) {
    group(id: $groupId) {
      id,
      name,
      upcomingEvents (input: {first: 1}) {
        edges {
          node {
            dateTime,
            timezone,
            shortUrl,
            tickets {
              edges {
                node {
                  user {
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
  }`;
  
  try {
    const response = await fetch('https://api.meetup.com/gql', {
      method: 'post',
      body: JSON.stringify({query, variables}),
      headers: {'Content-Type': 'application/json'}
    });
    
    const data = await response.json();
    
    return data.data.group.upcomingEvents.edges[0].node;
  } catch (e) {
    console.error(e)
    return null
  }
  
}

export function formatEventMessage ({ shortUrl, tickets}) {
  let namesText = '';
  
  if (!tickets) {
    namesText = 'us';
  } else {
    const userNames = tickets.edges.map(ticket => ticket.node.user.name);
  
    namesText = joinArrayHumanReadable(userNames)
  }
  
  return `Join ${namesText} at our event this evening! ${shortUrl}`
}
