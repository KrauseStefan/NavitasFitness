import { module } from 'angular';
import { adminPageComponent } from './AdminPageComponent';
import { AdminUiGridConstants } from './AdminUiGridsConstants';
import { adminTransacrtionGridComponent } from './transactionsGrid/AdminTransactionGrid';
import { adminUsersGridComponent } from './Users/AdminUsersGrid';

export { adminRouterState } from './AdminPageComponent';

export const adminPageModule = module('AdminPage', ['ui.grid', 'ui.grid.selection'])
  .service('adminUiGridConstants', AdminUiGridConstants)
  .component('adminUsersGrid', adminUsersGridComponent)
  .component('adminTransacrtionGrid', adminTransacrtionGridComponent)
  .component('adminPage', adminPageComponent);
