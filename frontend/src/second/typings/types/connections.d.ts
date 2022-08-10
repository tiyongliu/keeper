interface IConnection {
  server: string
  engine: string
  sshModeL: string
  user: string
  password: string
  _id: string
}

export interface IConnectionStatus {
  name: connectionStatus
  message: string
}

export interface IActiveConnection extends IConnection, IConnectionStatus {
}

type connectionStatus = 'pending' | 'ok' | 'error'
