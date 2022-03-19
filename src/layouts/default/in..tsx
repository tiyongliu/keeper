import { FC, MouseEventHandler, useMemo } from 'react';
import { DragBox } from '../sytyled';

interface DragProps {
  navigationRef: null | HTMLDivElement;
}
const MAX_WIDTH = window.innerWidth * 0.4;
const MIN_WIDTH = 200;

const Drag: FC<DragProps> = ({ navigationRef }) => {
  //拖拽
  return useMemo(() => {
    let firstX = 0;
    let firstW = 0;
    let canVisit = true;

    const setWidht = (width: number) => {
      if (navigationRef !== null && width < MAX_WIDTH && width > MIN_WIDTH) {
        // eslint-disable-next-line no-param-reassign
        navigationRef.style.width = `${width}px`;
      }
    };
    const onMouseMove = (e) => {
      if (!canVisit) {
        return;
      }
      canVisit = false;
      const nowX = e.pageX - firstX;
      setWidht(firstW + nowX);
      setTimeout(() => {
        canVisit = true;
      }, 10);
    };
    const onMouseLive = () => {
      document.removeEventListener('mousemove', onMouseMove, false);
      document.removeEventListener('mouseup', onMouseLive, false);
    };
    const onMouseDown: MouseEventHandler<HTMLDivElement> = (e) => {
      firstX = e.pageX;
      firstW = navigationRef?.clientWidth || 0;
      document.addEventListener('mousemove', onMouseMove, false);
      document.addEventListener('mouseup', onMouseLive, false);
    };

    return <DragBox onMouseDown={onMouseDown} />;
  }, [navigationRef]);
};

export default Drag;
