import { module } from 'angular';
import { adminAccessOverrideComponent } from './AccessOverride/AdminAccessOverride';
import { AdminAccessOverrideRest } from './AccessOverride/AdminAccessOverrideRest';
import { adminPageComponent } from './AdminPageComponent';
import { AdminUiGridConstants } from './AdminUiGridsConstants';
import { adminTransacrtionGridComponent } from './transactionsGrid/AdminTransactionGrid';
import { adminUsersGridComponent } from './Users/AdminUsersGrid';

export { adminRouterState } from './AdminPageComponent';

export const adminPageModule = module('AdminPage', ['ui.grid', 'ui.grid.selection', 'ui.grid.edit', 'ui.grid.cellNav'])
  .service('adminUiGridConstants', AdminUiGridConstants)
  .service('adminAccessOverrideRest', AdminAccessOverrideRest)
  .component('adminAccessOverride', adminAccessOverrideComponent)
  .component('adminUsersGrid', adminUsersGridComponent)
  .component('adminTransacrtionGrid', adminTransacrtionGridComponent)
  .component('adminPage', adminPageComponent);
