import registerCommand from './registerCommand'

registerCommand({
  id: 'new.connection',
  toolbar: true,
  icon: 'icon new-connection',
  toolbarName: 'Add connection',
  category: 'New',
  toolbarOrder: 1,
  name: 'Connection',
  testEnabled: () => false,
  onClick: () => {

  }
})
