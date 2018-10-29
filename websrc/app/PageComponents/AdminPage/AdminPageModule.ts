import { module } from 'angular';
import { adminAccessOverrideComponent } from './AccessOverride/AdminAccessOverride';
import { adminPageComponent } from './AdminPageComponent';
import { AdminUiGridConstants } from './AdminUiGridsConstants';
import { adminTransacrtionGridComponent } from './transactionsGrid/AdminTransactionGrid';
import { adminUsersGridComponent } from './Users/AdminUsersGrid';

export { adminRouterState } from './AdminPageComponent';

export const adminPageModule = module('AdminPage', ['ui.grid', 'ui.grid.selection', 'ui.grid.edit', 'ui.grid.cellNav'])
  .service('adminUiGridConstants', AdminUiGridConstants)
  .component('adminAccessOverride', adminAccessOverrideComponent)
  .component('adminUsersGrid', adminUsersGridComponent)
  .component('adminTransacrtionGrid', adminTransacrtionGridComponent)
  .component('adminPage', adminPageComponent);
