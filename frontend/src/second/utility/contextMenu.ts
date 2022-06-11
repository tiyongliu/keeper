
export function handleContextMenu(e, items: any = []) {
  e.preventDefault()
  e.stopPropagation()

  console.log(items, `items`)

  if (items) {
    const left = e.pageX;
    const top = e.pageY;
    console.log({left, top, items, targetElement: e.target})
  }

  if (items === '__no_menu') return

}
