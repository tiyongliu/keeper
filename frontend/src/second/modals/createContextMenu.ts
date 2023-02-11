import contextMenuVue from '/@/second/modals/ContextMenu.vue'
import { isClient } from '/@/utils/is';
import { CreateContextOptions, ContextMenuProps } from './typing';
import { createVNode, render } from 'vue';
import invalidateCommands from '/@/second/commands/invalidateCommands'
const menuManager: {
  domList: Element[];
  resolve: Fn;
} = {
  domList: [],
  resolve: () => {},
};

export const createContextMenu = function (options: CreateContextOptions) {
  const {event} = options || {};

  event && event?.preventDefault();

  if (!isClient) {
    return;
  }
  return new Promise(async (resolve) => {
    await invalidateCommands()

    const body = document.body;

    const container = document.createElement('div');
    const propsData: Partial<ContextMenuProps> = {};
    if (options.styles) {
      propsData.styles = options.styles;
    }

    if (options.items) {
      propsData.items = options.items;

      const left = event.pageX;
      const top = event.pageY;
      propsData.left = left
      propsData.top = top
    }

    if (options.event) {
      // @ts-ignore
      propsData.targetElement = event.target;
    }

    // @ts-ignore
    const vm = createVNode(contextMenuVue, {
      ...propsData,
      onClose: () => {
        //废弃
        // bootstrap.setCurrentDropDownMenu(null)
      }
    });
    render(vm, container);

    const handleClick = function () {
      menuManager.resolve('');
    };

    menuManager.domList.push(container);

    const remove = function () {
      menuManager.domList.forEach((dom: Element) => {
        try {
          dom && body.removeChild(dom);
        } catch (error) {}
      });
      body.removeEventListener('click', handleClick);
      body.removeEventListener('scroll', handleClick);
    };

    menuManager.resolve = function (arg) {
      remove();
      resolve(arg);
    };
    remove();
    body.appendChild(container);
    body.addEventListener('click', handleClick);
    body.addEventListener('scroll', handleClick);
  });
};

export const destroyContextMenu = function () {
  if (menuManager) {
    menuManager.resolve('');
    menuManager.domList = [];
  }
};
