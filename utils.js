function joinArrayHumanReadable (array) {
  let text = '';
  
  for (let i = 0; i < array.length; i++) {
    if (i === 0) {
      text += array[i];
      continue
    }
    
    if (array.length === 2) {
      text += ` and ${array[i]}`
      continue
    }
    
    if (i === array.length - 1) {
      text += `, and ${array[i]}`
      continue
    }
    
    text += `, ${array[i]}`;
  }
  
  return text;
}

module.exports = { joinArrayHumanReadable }
