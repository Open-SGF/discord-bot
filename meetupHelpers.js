import { joinArrayHumanReadable } from './utils.js';
import fetch from 'node-fetch';

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
    const response = await fetch.default('https://api.meetup.com/gql', {
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
