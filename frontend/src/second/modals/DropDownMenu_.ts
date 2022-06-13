function getElementOffset(element, side = null) {
  const de = document.documentElement;
  const box = element.getBoundingClientRect();
  const top = box.top + window.pageYOffset - de.clientTop;
  const left = box.left + window.pageXOffset - de.clientLeft;
  if (side == 'right') return {top: top, left: left + box.width};
  return {top: top, left: left};
}

export function fixPopupPlacement(element: HTMLElement): HTMLElement {
  const {width, height} = element.getBoundingClientRect()
  let offset = getElementOffset(element);

  let newLeft: null | number = null
  let newTop: null | number = null

  if (offset.left + width > window.innerWidth) {
    newLeft = offset.left - width

    if (newLeft < 0) newLeft = 0
  }

  if (offset.top + height > window.innerHeight) {
    newTop = offset.top - height

    if (newTop < 0) newTop = 0
    if (newTop + height > window.innerHeight) {
      element.style.height = `${window.innerHeight - newTop}px`
    }
  }

  if (newLeft != null) element.style.left = `${newLeft}px`
  if (newTop != null) element.style.top = `${newTop}px`

  return element
}
