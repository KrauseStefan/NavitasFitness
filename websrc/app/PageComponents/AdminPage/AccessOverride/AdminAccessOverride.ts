import { copy } from 'angular';
import 'angular-ui-grid';
import * as moment from 'moment';
import { AdminUiGridConstants } from '../AdminUiGridsConstants';
import { AccessIdOverrideDto, AdminAccessOverrideRest } from './AdminAccessOverrideRest';

interface AccessOverrideRow {
  dto: AccessIdOverrideDto;
  creationiDateStr: string;
  prevAccessId: string;
}

type GridOptions = uiGrid.IGridOptionsOf<AccessOverrideRow>;

export class AdminAccessOverrideCtrl {

  public readonly gridOptions: GridOptions;

  public accessOverrides: AccessOverrideRow[] = [];

  private readonly emptyRow: Readonly<AccessOverrideRow> = {
    dto: {
      accessId: '',
      startDate: moment(),
    },
    creationiDateStr: '-',
    prevAccessId: '',
  };

  private columnDefs: Array<uiGrid.IColumnDefOf<AccessOverrideRow>> = [
    { name: 'Access Id', field: 'dto.accessId' },
    { name: 'Creation Date', field: 'creationiDateStr', allowCellFocus: false, enableCellEdit: false },
    {
      field: 'delete',
      enableFiltering: false,
      allowCellFocus: false,
      name: '',
      width: 52,
      cellTemplate: `
        <md-button class="md-icon-button" ng-click="grid.appScope.$ctrl.rowAction(row)" aria-label="row action">
          <md-icon ng-if="row.entity.prevAccessId && row.entity.prevAccessId == row.entity.dto.accessId">
            delete
          </md-icon>
          <md-icon ng-if="row.entity.prevAccessId != row.entity.dto.accessId">
            save
          </md-icon>
        </md-button>`,
    },
  ];

  private readonly gridApi: ng.IPromise<uiGrid.IGridApiOf<AccessOverrideRow>>;

  private readonly specificGridOptions: GridOptions = {
    enableMinHeightCheck: false,
    enableCellEditOnFocus: true,
    modifierKeysToMultiSelectCells: true,
    data: this.accessOverrides,
    columnDefs: this.columnDefs,
  };

  constructor(
    adminUiGridConstants: AdminUiGridConstants,
    private uiGridConstants: uiGrid.IUiGridConstants,
    $q: ng.IQService,
    private adminAccessOverrideRest: AdminAccessOverrideRest,
  ) {
    const options = adminUiGridConstants.options;

    this.gridOptions = Object.assign({}, options, this.specificGridOptions);
    this.gridApi = $q((resolve) => this.gridOptions.onRegisterApi = resolve);

    adminAccessOverrideRest.getAllAccessIdOverrrides()
      .then((accessOverrides) => accessOverrides.map((ao) => {
        return {
          dto: ao,
          creationiDateStr: ao.startDate.format('DD-MM-YYYY'),
          prevAccessId: ao.accessId,
        };
      }))
      .then((accessOverrides) => {
        this.accessOverrides = accessOverrides;
        this.gridOptions.data = this.accessOverrides;
        this.addEmptyRow();
        this.updateRowsToShow();
      });
  }

  public rowAction(row: uiGrid.IGridRowOf<AccessOverrideRow>): void {
    const entity = row.entity;
    if (entity.prevAccessId !== entity.dto.accessId) {
      this.save(entity);
    } else {
      this.delete(entity);
      this.updateRowsToShow();
    }
  }

  private addEmptyRow(): void {
    const entity = copy(this.emptyRow);
    entity.dto.startDate = moment();

    this.accessOverrides.unshift(entity);
  }

  private save(entity: AccessOverrideRow) {
    entity.dto.startDate = moment();

    this.adminAccessOverrideRest.saveAccessIdOverride(entity.dto).then(() => {
      if (entity.prevAccessId === '') {
        this.addEmptyRow();
        this.updateRowsToShow();
      }
      entity.prevAccessId = entity.dto.accessId;
      entity.creationiDateStr = entity.dto.startDate.format('DD-MM-YYYY');

    });
  }

  private delete(entity: AccessOverrideRow) {
    const index = this.accessOverrides.indexOf(entity);
    if (index !== -1) {
      this.adminAccessOverrideRest.deleteAccessIdOverride(entity.dto.accessId)
        .then(() => {
          this.accessOverrides.splice(index, 1);
          this.updateRowsToShow();
        });
    } else {
      throw new Error('Cannot delete row, row not found');
    }
  }

  private updateRowsToShow() {

    this.gridApi.then((api) => {
      api.core.notifyDataChange(this.uiGridConstants.dataChange.ROW);
    });
  }

}

export const adminAccessOverrideComponent: ng.IComponentOptions = {
  controller: AdminAccessOverrideCtrl,
  templateUrl: '/PageComponents/AdminPage/AccessOverride/AdminAccessOverride.html',
};
