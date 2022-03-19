export function safeJsonParse(json, defaultValue?, logError = false) {
  console.log(json, defaultValue, logError)
  try {
    return JSON.parse(json);
  } catch (err) {
    if (logError) {
      console.error(`Error parsing JSON value "${json}"`, err);
    }
    return defaultValue;
  }
}
