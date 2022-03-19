import {EngineDriver, ExtensionsDirectory} from '/@/second/keeper-types'
import {isPlainObject, isString} from 'lodash-es'

export function findEngineDriver(connection, extensions: ExtensionsDirectory): EngineDriver {
  if (isString(connection)) {
    return extensions.drivers.find(x => x.engine == connection)!
  }

  if (isPlainObject(connection)) {
    const {engine} = connection
    if (engine) {
      return extensions.drivers.find(x => x.engine == engine)!
    }
  }

  return null!
}
