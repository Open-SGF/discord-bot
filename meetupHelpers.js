const { joinArrayHumanReadable } = require('./utils');

async function fetchNextEventForGroup (groupId) {
  const fetch = await import('node-fetch');
  
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
  
  const response = await fetch.default('https://api.meetup.com/gql', {
    method: 'post',
    body: JSON.stringify({query, variables}),
    headers: {'Content-Type': 'application/json'}
  });
  
  const data = await response.json();
  
  return data.data.group.upcomingEvents.edges[0].node;
}

function formatEventMessage ({ shortUrl, tickets}) {
  const userNames = tickets.edges.map(ticket => ticket.node.user.name);
  
  const namesText = joinArrayHumanReadable(userNames)
  
  return `Join ${namesText} at our event this evening! ${shortUrl}`
}

module.exports = { fetchNextEventForGroup, formatEventMessage }
