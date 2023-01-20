/**
 * Configure and register global directives
 */
import type { App } from 'vue';
import { setupPermissionDirective } from './permission';
import { setupLoadingDirective } from './loading';
import { setupSplitterDrag } from './splitterDrag';
import { setupResizeObserver } from './resizeObserver';

export function setupGlobDirectives(app: App) {
  setupPermissionDirective(app);
  setupLoadingDirective(app);
  setupSplitterDrag(app);
  setupResizeObserver(app);
}
