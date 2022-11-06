import {inject, onBeforeUnmount} from 'vue'
import {padStart,} from 'lodash-es'
import {updateStatusBarInfoItem} from '/@/second/utility/statusBarStore'

function formatSeconds(duration) {
  if (duration == null) return '';
  const hours = padStart(Math.floor(duration / 3600).toString(), 2, '0');
  const minutes = padStart((Math.floor(duration / 60) % 60).toString(), 2, '0');
  const seconds = padStart((duration % 60).toString(), 2, '0');
  return `${hours}:${minutes}:${seconds}`;
}

export default function useTimerLabel() {
  let duration: number | null = null
  let timerHandle: number | null = null
  const tabid = inject('tabid')

  const update = () => {
    updateStatusBarInfoItem(tabid, 'durationSeconds', {text: formatSeconds(duration)})
  }

  const start = () => {
    duration = 0;
    update();
    timerHandle = window.setInterval(() => {
      duration! += 1;
      update();
    }, 1000);
  }

  const stop = () => {
    update();
    if (timerHandle) {
      window.clearInterval(timerHandle);
      timerHandle = null;
    }
  }

  onBeforeUnmount(() => {
    if (timerHandle) {
      window.clearInterval(timerHandle);
      timerHandle = null
    }
  })

  return {
    start,
    stop,
  }
}
