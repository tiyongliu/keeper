export interface IConnectionAppObjectData {
  server: string
  engine: string
  sshMode: string
  sshPort: string
  sshKeyfile: string
  port: string
  user: string
  password: string
  _id: string
  status: {
    name: string
    message?: string
  }

  singleDatabase?: string
  defaultDatabase: string
  isReadOnly?: boolean
}

export interface IConnections {
  engine: string
  password: string
  server: string
  sshKeyfile: string
  sshMode: string
  sshPort: string
  user: string
  _id: string
  singleDatabase?: string
}

export interface IPinnedDatabasesItem {
  extInfo?: string | Ref<string> | undefined;
  connection: IConnections
  name: string
  objectTypeField?: string

  title: string
}

type IIsExpandable = (data: { _id: string, singleDatabase: boolean }) => boolean

export interface IPinnedTablesItem {
  pureName: string
  tableRowCount: string
  modifyDate: string
  objectId: string
  contentHash: string
  columns: { [key in string] }[]
  primaryKey: {
    constraintName: string
    pureName: string
    constraintType: string
    columns: { columnName: string }[]
  }
  foreignKeys: {
    constraintName: string
    constraintType: string
    pureName: string
    refTableName: string
    updateAction: string
    deleteAction: string
  }[]
}
