import {ApplicationDefinition} from '/@/second/keeper-types'

export function filterAppsForDatabase(connection, database: string, $apps): ApplicationDefinition[] {
  const db = (connection?.databases || []).find(x => x.name == database);
  return $apps?.filter(app => db && db[`useApp:${app.name}`]);
}
