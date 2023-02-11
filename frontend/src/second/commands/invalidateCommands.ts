import {useBootstrapStoreWithOut} from "/@/store/modules/bootstrap";

export interface SubCommand {
  text: string;
  onClick: Function;
}

export interface GlobalCommand {
  id: string;
  category: string; // null for group commands
  isGroupCommand?: boolean;
  name: string;
  text?: string /* category: name */;
  keyText?: string;
  keyTextFromGroup?: string; // automatically filled from group
  group?: string;
  getSubCommands?: () => SubCommand[];
  onClick?: Function;
  testEnabled?: () => boolean;
  // enabledStore?: any;
  icon?: string;
  toolbar?: boolean;
  enabled?: boolean;
  showDisabled?: boolean;
  toolbarName?: string;
  menuName?: string;
  toolbarOrder?: number;
  disableHandleKeyText?: string;
  isRelatedToTab?: boolean,
  systemCommand?: boolean;
}

let isInvalidated = false

export default async function invalidateCommands() {
  const bootstrap = useBootstrapStoreWithOut()
  if (isInvalidated) return
  isInvalidated = true

  isInvalidated = false

  bootstrap.updateCommands(dct => {
    let res: any = null;
    for (const command of Object.values(dct) as GlobalCommand[]) {
      if (command.isGroupCommand) continue;
      const { testEnabled } = command;
      let enabled = command.enabled;
      if (testEnabled) enabled = testEnabled();
      if (enabled != command.enabled) {
        if (!res) res = { ...dct };
        res[command.id].enabled = enabled;
      }
    }
    if (res) {
      const values = Object.values(res) as GlobalCommand[];
      // test enabled for group commands
      for (const command of values) {
        if (!command.isGroupCommand) continue;
        const groupSources = values.filter(x => x.group == command.group && !x.isGroupCommand && x.enabled);
        command.enabled = groupSources.length > 0;
        // for (const source of groupSources) {
        //   source.keyTextFromGroup = command.keyText;
        // }
      }
    }
    return res || dct;
  })
}
let isInvalidatedDefinitions = false;

export async function invalidateCommandDefinitions(bootstrap) {
  if (isInvalidatedDefinitions) return;
  isInvalidatedDefinitions = true;

  bootstrap.updateCommands(dct => {
    const res = { ...dct };
    const values = Object.values(res) as GlobalCommand[];
    // test enabled for group commands
    for (const command of values) {
      if (!command.isGroupCommand) continue;
      const groupSources = values.filter(x => x.group == command.group && !x.isGroupCommand);

      for (const source of groupSources) {
        source.keyTextFromGroup = command.keyText;
      }
    }
    return res;
  })
}
