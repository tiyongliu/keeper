// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export function ListDatabases(arg1) {
  return window['go']['bridge']['ServerConnections']['ListDatabases'](arg1);
}

export function Ping(arg1) {
  return window['go']['bridge']['ServerConnections']['Ping'](arg1);
}

export function Refresh(arg1) {
  return window['go']['bridge']['ServerConnections']['Refresh'](arg1);
}

export function ServerStatus() {
  return window['go']['bridge']['ServerConnections']['ServerStatus']();
}

export function Close(arg1, arg2) {
  window['go']['bridge']['ServerConnections']['Close'](arg1, arg2);
}
