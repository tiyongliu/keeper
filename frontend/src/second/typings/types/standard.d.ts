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

/*
*
connection: {engine: "mysql@dbgate-pluginMysql",…}
engine: "mysql@dbgate-pluginMysql"
password: "crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA=="
server: "localhost"
sshKeyfile: "C:\\Users\\Administrator\\.ssh\\id_rsa"
sshMode: "userPassword"
sshPort: "22"
user: "root"
_id: "065caa90-a8c6-11ec-9b4b-6f98950c4d7a"
name: "mysql"
* */
