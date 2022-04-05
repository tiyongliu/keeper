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

  defaultDatabase?: string

  isReadOnly?: boolean
}
