import _ from 'lodash';

export function sortTabs(tabs: any[]): any[]{

  // _.sortBy(collection, [iteratees=[_.identity]])
  // 创建一个元素数组。 以 iteratee 处理的结果升序排序。 这个方法执行稳定排序，也就是说相同元素会保持原始排序。 iteratees 调用1个参数： (value)
  return _.sortBy(tabs, [ x => x.tabOrder || 0, x => getTabDbKey(x),'title','tabid'])
}

export function getTabDbKey(tab){
  if(tab.props && tab.props.conid && tab.props.database){
    return `database://${tab.props.data}-${tab.props.conid}`;
  }
  if (tab.props && tab.props.conid) {
    return `server://${tab.props.conid}`;
  }
  if (tab.props && tab.props.archiveFolder) {
    return `archive://${tab.props.archiveFolder}`;
  }
  return null;
}

export function groupTabs(tabs: any[]){
  const res = [];
  for( const tab of sortTabs(tabs)){

  }
}
