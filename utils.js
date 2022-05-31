function joinArrayHumanReadable (array) {
  if (array.length <= 2) { return array.join(' and '); }
  const rest = array.pop();
  return array.join(', ') + ', and ' + rest;
}

module.exports = { joinArrayHumanReadable }
