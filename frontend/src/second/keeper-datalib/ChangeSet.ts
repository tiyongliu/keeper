import {isEqual, maxBy} from 'lodash-es'
import {NamedObjectInfo} from '/@/second/keeper-types'

export interface ChangeSetItem {
  pureName: string;
  schemaName?: string;
  insertedRowIndex?: number;
  document?: any;
  condition?: { [column: string]: string };
  fields?: { [column: string]: string };
}

export interface ChangeSet {
  inserts: ChangeSetItem[];
  updates: ChangeSetItem[];
  deletes: ChangeSetItem[];
}

export interface ChangeSetRowDefinition {
  pureName: string;
  schemaName: string;
  insertedRowIndex?: number;
  condition?: { [column: string]: string };
}

export interface ChangeSetFieldDefinition extends ChangeSetRowDefinition {
  uniqueName: string;
  columnName: string;
}

export function createChangeSet(): ChangeSet {
  return {
    inserts: [],
    updates: [],
    deletes: [],
  };
}

export function revertChangeSetRowChanges(changeSet: ChangeSet, definition: ChangeSetRowDefinition): ChangeSet {
  // console.log('definition', definition);
  const [field, item] = findExistingChangeSetItem(changeSet, definition);
  // console.log('field, item', field, item);
  // console.log('changeSet[field]', changeSet[field]);
  // console.log('changeSet[field] filtered', changeSet[field].filter((x) => x != item);
  if (item)
    return {
      ...changeSet,
      [field]: changeSet[field].filter(x => x != item),
    };
  return changeSet;
}

export function findExistingChangeSetItem(
  changeSet: ChangeSet,
  definition: ChangeSetRowDefinition
): [keyof ChangeSet, Nullable<ChangeSetItem>] {
  if (!changeSet || !definition) return ['updates', null];
  if (definition.insertedRowIndex != null) {
    return [
      'inserts',
      changeSet.inserts.find(
        x =>
          x.pureName == definition.pureName &&
          x.schemaName == definition.schemaName &&
          x.insertedRowIndex == definition.insertedRowIndex
      )!,
    ];
  } else {
    const inUpdates = changeSet.updates.find(
      x =>
        x.pureName == definition.pureName &&
        x.schemaName == definition.schemaName &&
        isEqual(x.condition, definition.condition)
    );
    if (inUpdates) return ['updates', inUpdates];

    const inDeletes = changeSet.deletes.find(
      x =>
        x.pureName == definition.pureName &&
        x.schemaName == definition.schemaName &&
        isEqual(x.condition, definition.condition)
    );
    if (inDeletes) return ['deletes', inDeletes];

    return ['updates', null];
  }
}

export function setChangeSetValue(
  changeSet: ChangeSet,
  definition: ChangeSetFieldDefinition,
  value: string
): ChangeSet {
  if (!changeSet || !definition) return changeSet;
  let [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  if (fieldName == 'deletes') {
    changeSet = revertChangeSetRowChanges(changeSet, definition);
    [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  }
  if (existingItem) {
    return {
      ...changeSet,
      [fieldName]: changeSet[fieldName].map(item =>
        item == existingItem
          ? {
            ...item,
            fields: {
              ...item.fields,
              [definition.uniqueName]: value,
            },
          }
          : item
      ),
    };
  }

  return {
    ...changeSet,
    [fieldName]: [
      ...changeSet[fieldName],
      {
        pureName: definition.pureName,
        schemaName: definition.schemaName,
        condition: definition.condition,
        insertedRowIndex: definition.insertedRowIndex,
        fields: {
          [definition.uniqueName]: value,
        },
      },
    ],
  };
}

export function setChangeSetRowData(
  changeSet: ChangeSet,
  definition: ChangeSetRowDefinition,
  document: any
): ChangeSet {
  if (!changeSet || !definition) return changeSet;
  let [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  if (fieldName == 'deletes') {
    changeSet = revertChangeSetRowChanges(changeSet, definition);
    [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  }
  if (existingItem) {
    return {
      ...changeSet,
      [fieldName]: changeSet[fieldName].map(item =>
        item == existingItem
          ? {
            ...item,
            fields: {},
            document,
          }
          : item
      ),
    };
  }

  return {
    ...changeSet,
    [fieldName]: [
      ...changeSet[fieldName],
      {
        pureName: definition.pureName,
        schemaName: definition.schemaName,
        condition: definition.condition,
        insertedRowIndex: definition.insertedRowIndex,
        document,
      },
    ],
  };
}

export function getChangeSetInsertedRows(changeSet: ChangeSet, name?: NamedObjectInfo) {
  if (!name) return []
  if (!changeSet) return []
  const rows = changeSet.inserts.filter(x => x.pureName == name.pureName && x.schemaName == name.schemaName)
  const maxIndex = maxBy(rows, x => x.insertedRowIndex)?.insertedRowIndex
  if (maxIndex == null) return []
  const res = Array(maxIndex + 1).fill({})
  for (const row of rows) {
    res[row.insertedRowIndex!] = row.fields
  }
  return res
}

function consolidateInsertIndexes(changeSet: ChangeSet, name: NamedObjectInfo): ChangeSet {
  const indexes = changeSet.inserts
    .filter(x => x.pureName == name.pureName && x.schemaName == name.schemaName)
    .map(x => x.insertedRowIndex);

  indexes.sort((a, b) => a! - b!);
  if (indexes[indexes.length - 1] != indexes.length - 1) {
    return {
      ...changeSet,
      inserts: changeSet.inserts.map(x => ({
        ...x,
        insertedRowIndex: indexes.indexOf(x.insertedRowIndex),
      })),
    };
  }

  return changeSet;
}

export function deleteChangeSetRows(changeSet: ChangeSet, definition: ChangeSetRowDefinition): ChangeSet {
  let [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  if (fieldName == 'updates') {
    changeSet = revertChangeSetRowChanges(changeSet, definition);
    [fieldName, existingItem] = findExistingChangeSetItem(changeSet, definition);
  }
  if (fieldName == 'inserts') {
    return consolidateInsertIndexes(revertChangeSetRowChanges(changeSet, definition), definition);
  } else {
    if (existingItem && fieldName == 'deletes') return changeSet;
    return {
      ...changeSet,
      deletes: [
        ...changeSet.deletes,
        {
          pureName: definition.pureName,
          schemaName: definition.schemaName,
          condition: definition.condition,
        },
      ],
    };
  }
}

export function changeSetInsertNewRow(changeSet: ChangeSet, name?: NamedObjectInfo): ChangeSet {
  // console.log('INSERT', name);
  const insertedRows = getChangeSetInsertedRows(changeSet, name);
  return {
    ...changeSet,
    inserts: [
      ...changeSet.inserts,
      {
        ...name,
        insertedRowIndex: insertedRows.length,
        fields: {},
      },
    ],
  } as ChangeSet;
}

export function changeSetInsertDocuments(changeSet: ChangeSet, documents: any[], name?: NamedObjectInfo): ChangeSet {
  const insertedRows = getChangeSetInsertedRows(changeSet, name);
  return {
    ...changeSet,
    inserts: [
      ...changeSet.inserts,
      ...documents.map((doc, index) => ({
        ...name,
        insertedRowIndex: insertedRows.length + index,
        fields: doc,
      })),
    ],
  } as ChangeSet;
}

export function changeSetContainsChanges(changeSet: ChangeSet) {
  if (!changeSet) return false;
  return changeSet.deletes.length > 0 || changeSet.updates.length > 0 || changeSet.inserts.length > 0;
}
