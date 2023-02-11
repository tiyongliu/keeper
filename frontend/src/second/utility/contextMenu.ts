import {inject, provide, Ref, unref} from 'vue'
import {compact, intersection, isArray, isFunction, isPlainObject} from 'lodash-es'
import getAsArray from './getAsArray'
import {ContextMenuItem} from "/@/second/modals/typing";
import invalidateCommands from '/@/second/commands/invalidateCommands'

export function registerMenu(originMenu, ...items) {
  const parentMenu = inject('componentContextMenu') as any
  originMenu.push(...[unref(parentMenu), ...unref(items)])
  provide('componentContextMenu', originMenu)
}

export function getContextMenu() {
  return inject('componentContextMenu') as Ref<ContextMenuItem[]>
}

export async function handleContextMenu(e, items: any = []) {
  e.preventDefault()
  e.stopPropagation()

  await invalidateCommands()

  if (items) {
    // const left = e.pageX
    // const top = e.pageY
    // bootstrap.setCurrentDropDownMenu({left, top, items, targetElement: e.target})
  }

  if (items === '__no_menu') return
}

function doExtractMenuItems(menu, res, options) {
  if (isFunction(menu)) {
    doExtractMenuItems(menu(options), res, options);
  } else if (isArray(menu)) {
    for (const item of menu) {
      doExtractMenuItems(item, res, options);
    }
  } else if (isPlainObject(menu) && !menu._skip) {
    res.push(menu);
  }
}

function extractMenuItems(menu, options = null) {
  let res = [];
  doExtractMenuItems(menu, res, options);
  res = processTags(res);
  return res;
}

function processTags(items) {
  const res = [];
  const tagged = [];

  for (const menu of items.filter(x => x.tag)) {
    // @ts-ignore
    tagged.push({...menu, tags: getAsArray(menu.tag)});
  }

  for (const menu of items.filter(x => !x.tag)) {
    if (menu.placeTag) {
      const placeTags = getAsArray(menu.placeTag);
      for (let index = 0; index < tagged.length;) {
        const current = tagged[index];
        // @ts-ignore
        if (intersection(placeTags, current.tags).length > 0) {
          tagged.splice(index, 1);
          res.push(current);
        } else {
          index++;
        }
      }
    } else {
      // @ts-ignore
      res.push(menu);
    }
  }

  // append tagged, which were not appended by placeTag
  res.push(...tagged);

  return res
}

function filterMenuItems(items) {
  const res = [];
  let wasDivider = false;
  let wasItem = false;
  for (const item of items.filter(x => !x.disabled || !x.hideDisabled)) {
    if (item.divider) {
      if (wasItem) {
        wasDivider = true;
      }
    } else {
      if (wasDivider) {
        // @ts-ignore
        res.push({divider: true});
      }
      wasDivider = false;
      wasItem = true;
      // @ts-ignore
      res.push(item);
    }
  }
  return res;
}

function mapItem(item, bootstrap) {
  const commands = bootstrap.getCommandsCustomized
  if (item.command) {
    const command = commands[item.command];
    if (command) {
      return {
        text: item.text || command.menuName || command.toolbarName || command.name,
        keyText: command.keyText || command.keyTextFromGroup,
        onClick: () => {
          if (command.getSubCommands) bootstrap.setVisibleCommandPalette(command)
          else if (command.onClick) command.onClick();
        },
        disabled: !command.enabled,
        hideDisabled: item.hideDisabled,
      }
    }
    return null
  }
  return item
}

export function prepareMenuItems(items, options, bootstrap) {
  const extracted = extractMenuItems(items, options);
  const compacted = compact(extracted.map(x => mapItem(x, bootstrap)))
  return filterMenuItems(compacted);
}
