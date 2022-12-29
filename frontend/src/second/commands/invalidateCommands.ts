let isInvalidated = false
export default async function invalidateCommands() {
  if (isInvalidated) return
  isInvalidated = true

  isInvalidated = false
}
