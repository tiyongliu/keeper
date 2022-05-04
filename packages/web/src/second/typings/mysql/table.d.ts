type TableName = { name: string }

export interface TablesNameSort extends TableName {
  sortOrder?: string
}

export interface basisTablesField {

}

interface PureName {
  pureName: string
  tableRowCount: string
  modifyDate: Date.string
}

// ["pureName", "tableRowCount", "modifyDate", "objectId", "contentHash"]


type Column = Record<>

  /*
  *
 type basisConnField = Record<keyof ['server',
  'engine', 'sshMode', 'sshPort', 'sshKeyfile',
  'user', 'password'
], string>
  * */
