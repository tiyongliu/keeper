import {TablesNameSort} from './table'
type ConnStatus = 'ok' | 'error' | 'pending'

type basisConnField = Record<keyof ['server',
  'engine', 'sshMode', 'sshPort', 'sshKeyfile',
  'user', 'password'
], string>

export interface ConnectionsWithStatus extends basisConnField {
  readonly _id: string
  status?: {
    name: ConnStatus
    message?: string
  }

  singleDatabase?: string
  defaultDatabase: string
  isReadOnly?: boolean
  engine?: string
  port?: string
}


