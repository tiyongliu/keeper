export function getDatabaseFileLabel(databaseFile) {
  if (!databaseFile) return databaseFile;
  const m = databaseFile.match(/[\/]([^\/]+)$/);
  if (m) return m[1];
  return databaseFile;
}

function getConnectionLabelCore(connection, { allowExplicitDatabase = true } = {}) {
  if (!connection) {
    console.log(`gjkldgjlfdgfdgdggdgggg1111111111111111111111111`)
    return null;
  }
  if (connection.displayName) {
    return connection.displayName;
  }
  if (connection.singleDatabase && connection.host && allowExplicitDatabase && connection.defaultDatabase) {
    return `${connection.defaultDatabase} on ${connection.host}`;
  }
  if (connection.databaseFile) {
    return getDatabaseFileLabel(connection.databaseFile);
  }
  if (connection.host) {
    return connection.host;
  }
  if (connection.singleDatabase && connection.defaultDatabase) {
    return `${connection.defaultDatabase}`;
  }

  return '';
}

export default function getConnectionLabel(connection, { allowExplicitDatabase = true, showUnsaved = false } = {}) {
  const res = getConnectionLabelCore(connection, { allowExplicitDatabase });

  if (res && showUnsaved && connection?.unsaved) {
    return `${res} - unsaved`;
  }

  return res;
}
