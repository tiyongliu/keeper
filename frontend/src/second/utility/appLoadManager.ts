let appIsLoaded = false;
let onLoad: Function[] = [];

export function setAppLoaded() {
  appIsLoaded = true;
  for (const func of onLoad) {
    func();
  }
  onLoad = [];
}

export function getAppLoaded() {
  return appIsLoaded;
}

export function callWhenAppLoaded(callback: Function) {
  if (appIsLoaded) {
    callback();
  } else {
    onLoad.push(callback);
  }
}
