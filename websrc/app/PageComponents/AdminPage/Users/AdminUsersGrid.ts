import { copy } from 'angular';
import 'angular-ui-grid';
// tslint:disable-next-line no-implicit-dependencies
import { selection } from 'ui-grid';
import { AdminUiGridConstants } from '../AdminUiGridsConstants';
import { User } from './User';

type IUserGridOptions = selection.IGridOptions & uiGrid.IGridOptionsOf<User> & { data: User[] };

export class AdminUsersGridCtrl {

  public usersBackup: User[] = [];
  public users: User[] = [];
  public readonly gridOptionsUsers: IUserGridOptions;

  private gridApiDefered = this.$q.defer<uiGrid.IGridApiOf<User>>();

  private readonly gridOptionsUsersSpecific: uiGrid.IGridOptions = {
    enableRowSelection: true,
    modifierKeysToMultiSelect: true,
    noUnselect: false,
    multiSelect: true,
    onRegisterApi: (gridApi) => {
      this.gridApiDefered.resolve(gridApi);
      this.configureSelectionListeners(gridApi);
    },
  };

  constructor(
    private $scope: ng.IScope,
    private $q: ng.IQService,
    $http: ng.IHttpService,
    private uiGridConstants: uiGrid.IUiGridConstants,
    adminUiGridConstants: AdminUiGridConstants,
  ) {
    const options = adminUiGridConstants.options;
    // small hack, I suspect headerRowHeight should have been public in the typings

    this.gridOptionsUsers = <any>Object.assign({}, options, this.gridOptionsUsersSpecific);

    $http.get<{ users: User[], keys: string[] }>('/rest/user/all').then((res) => {
      this.users = res.data.users;
      this.gridOptionsUsers.data = this.users;

      this.gridOptionsUsers.minRowsToShow = 10;
      if (this.users.length < this.gridOptionsUsers.minRowsToShow) {
        this.gridOptionsUsers.minRowsToShow = this.users.length;
        this.gridOptionsUsers.enableVerticalScrollbar = this.uiGridConstants.scrollbars.NEVER;
      }

      res.data.keys.forEach((key: string, i: number) => {
        this.users[i].key = key;
      });

      this.usersBackup = copy(this.users);
      this.sortBy('accessId');
    });
  }

  public makeUsersUnique(testValue: keyof User) {
    let prev: User = <any>{};
    this.users = this.users.reduce((acc, value) => {
      if (prev && prev[testValue] === value[testValue]) {
        if (acc.length > 0 && acc[acc.length - 1] !== prev) {
          acc.push(prev);
        }

        acc.push(value);
      }

      prev = value;
      return acc;
    }, <User[]>[]);

    return this.users;
  }

  public sortBy(testValue: keyof User) {
    this.users = this.users.sort((a, b) => a[testValue].localeCompare(b[testValue]));
    return this.users;
  }

  private configureSelectionListeners(gridApi: uiGrid.IGridApiOf<User>): void {
    const displayTransactions = () => {
      this.userSelectionUpdated({
        users: gridApi.selection.getSelectedRows(),
      });
    };

    gridApi.selection.on.rowSelectionChanged(
      this.$scope,
      () => displayTransactions(),
    );

    gridApi.selection.on.rowSelectionChangedBatch(
      this.$scope,
      () => displayTransactions(),
    );
  }

  private userSelectionUpdated: (values: { users: User[] }) => void = () => { throw new Error('Not implemented'); };

}

export const adminUsersGridComponent: ng.IComponentOptions = {
  bindings: {
    userSelectionUpdated: '&',
  },
  controller: AdminUsersGridCtrl,
  templateUrl: '/PageComponents/AdminPage/Users/AdminUsersGrid.html',
};
