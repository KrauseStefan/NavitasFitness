import 'angular-ui-grid';
import * as moment from 'moment';
// tslint:disable-next-line no-implicit-dependencies
// import { selection } from 'ui-grid';
import { AdminUiGridConstants } from '../AdminUiGridsConstants';

interface AccessOverrideRow {
  dirty?: boolean;
  id: string;
  startDate: moment.Moment | null; // Moment;
  endDate: moment.Moment | null; // Moment;'
}

type GridOptions = uiGrid.IGridOptionsOf<AccessOverrideRow>
  & { data: AccessOverrideRow[] };

export class AdminAccessOverrideCtrl {

  public readonly gridOptions: GridOptions;

  public accessOverrides: AccessOverrideRow[] = [
    { id: 'accessId1', startDate: moment(), endDate: moment().add(6, 'months'), dirty: true },
    { id: 'accessId2', startDate: moment(), endDate: moment().add(6, 'months') },
    { id: 'accessId3', startDate: moment(), endDate: moment().add(6, 'months') },
    { id: '', startDate: null, endDate: null, dirty: true },
  ];

  private columnDefs: Array<uiGrid.IColumnDefOf<AccessOverrideRow>> = [
    { name: 'Access Id', field: 'id' },
    { name: 'Start Date', field: 'startDate' },
    { name: 'End Date', field: 'endDate' },
    {
      field: 'delete',
      enableFiltering: false,
      allowCellFocus: false,
      name: '',
      width: 52,
      cellTemplate: `
        <md-button class="md-icon-button" ng-click="grid.appScope.$ctrl.rowAction(row)" aria-label="row action">
          <md-icon ng-if="!row.entity.dirty">delete</md-icon>
          <md-icon ng-if="row.entity.dirty">done</md-icon>
        </md-button>`,
    },
  ];

  private readonly gridApi: ng.IPromise<uiGrid.IGridApiOf<AccessOverrideRow>>;

  private readonly specificGridOptions: GridOptions = {
    enableCellEditOnFocus: true,
    modifierKeysToMultiSelectCells: true,
    data: this.accessOverrides,
    columnDefs: this.columnDefs,
  };

  constructor(
    adminUiGridConstants: AdminUiGridConstants,
    private uiGridConstants: uiGrid.IUiGridConstants,
    $q: ng.IQService,
    private $timeout: ng.ITimeoutService,
  ) {
    const options = adminUiGridConstants.options;

    this.gridOptions = Object.assign({}, options, this.specificGridOptions);
    this.gridApi = $q((resolve) => this.gridOptions.onRegisterApi = resolve);

    this.updateMinRowsToShow();
  }

  public rowAction(row: uiGrid.IGridRowOf<AccessOverrideRow>): void {
    const entity = row.entity;
    if (entity.dirty) {
      this.save(entity);
    } else {
      this.delete(entity);
      this.updateMinRowsToShow();
    }
  }

  private save(entity: AccessOverrideRow) {
    throw new Error('Save not implemented' + entity.id);
  }

  private delete(entity: AccessOverrideRow) {
    const index = this.accessOverrides.indexOf(entity);
    if (index !== -1) {
      this.accessOverrides.splice(index, 1);
      this.updateMinRowsToShow();
      throw new Error('Delete not implemented yet');
    } else {
      throw new Error('Cannot delete row, row not found');
    }
  }

  private updateMinRowsToShow() {
    this.gridOptions.minRowsToShow = this.gridOptions.minRowsToShow || 10;
    if (this.accessOverrides.length < this.gridOptions.minRowsToShow) {
      this.gridOptions.minRowsToShow = this.accessOverrides.length;
      this.gridOptions.enableVerticalScrollbar = this.uiGridConstants.scrollbars.NEVER;
      this.reCreateGrid();
    }
    this.gridApi.then((api) => {
      api.core.notifyDataChange(this.uiGridConstants.dataChange.ROW);
    });
  }

  private reCreateGrid() {
    const backup = this.accessOverrides;
    this.accessOverrides = [];
    this.$timeout(() => this.accessOverrides = backup);
  }
}

export const adminAccessOverrideComponent: ng.IComponentOptions = {
  controller: AdminAccessOverrideCtrl,
  templateUrl: '/PageComponents/AdminPage/AccessOverride/AdminAccessOverride.html',
};
