import { copy, IComponentOptions, material } from 'angular';
// tslint:disable-next-line no-implicit-dependencies
import { selection } from 'ui-grid';

type IGridOptions = selection.IGridOptions & uiGrid.IGridOptionsOf<IUser>;

interface IUser {
  name: string;
  email: string;
  accessId: string;
  key: string;
}

interface ITransaction {
  amount: number;
  currency: string;
  isActive: boolean;
  paymentDate: string;
  status: string;
}

export const adminRouterState: angular.ui.IState = {
  template: '<admin-page></admin-page>',
  url: '/admin/',
};

export class AdminPageCtrl {

  public usersBackup: IUser[] = [];
  public users: IUser[] = [];
  public transactions: ITransaction[] = [];

  public gridOptions: IGridOptions = {
    data: [],
    enableColumnMenus: false,
    enableFiltering: true,
    enableHorizontalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    enableVerticalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    enableRowHeaderSelection: false,
    enableRowSelection: true,
    modifierKeysToMultiSelect: true,
    multiSelect: true,
    noUnselect: false,
    onRegisterApi: (gridApi) => {
      this.gridApiDefered.resolve(gridApi);

      const displayTransactions = () => {
        this.displayTransactions(gridApi.selection.getSelectedRows());
      };

      gridApi.selection.on.rowSelectionChanged(
        this.$scope,
        () => displayTransactions(),
      );

      gridApi.selection.on.rowSelectionChangedBatch(
        this.$scope,
        () => displayTransactions(),
      );
    },
    rowHeight: 42,
  };

  private transactionsCache: { [key: string]: ITransaction[] } = {};
  private gridApiDefered = this.$q.defer<uiGrid.IGridApiOf<IUser>>();

  constructor(
    private $q: ng.IQService,
    private $scope: ng.IScope,
    private $http: ng.IHttpService,
    private uiGridConstants: uiGrid.IUiGridConstants,
    private $mdToast: material.IToastService,
  ) {
    const headerHeight = 33;
    const filterHeight = 28;

    $http.get<{ users: IUser[], keys: string[] }>('/rest/user/all').then((res) => {
      this.users = res.data.users;
      this.gridOptions.data = this.users;

      // small hack, I suspect headerRowHeight should have been public in the typings
      (<any>this.gridOptions).headerRowHeight = Math.ceil((filterHeight + headerHeight) / 2);

      this.gridOptions.minRowsToShow = 10;
      if (this.users.length < this.gridOptions.minRowsToShow) {
        this.gridOptions.minRowsToShow = this.users.length;
        this.gridOptions.enableVerticalScrollbar = this.uiGridConstants.scrollbars.NEVER;
      }

      res.data.keys.forEach((key: string, i: number) => {
        this.users[i].key = key;
      });

      this.usersBackup = copy(this.users);
      this.sortBy('accessId');
    });
  }

  public makeUsersUnique(testValue: keyof IUser) {
    let prev: IUser = <any>{};
    this.users = this.users.reduce((acc, value) => {
      if (prev && prev[testValue] === value[testValue]) {
        if (acc.length > 0 && acc[acc.length - 1] !== prev) {
          acc.push(prev);
        }

        acc.push(value);
      }

      prev = value;
      return acc;
    }, <IUser[]>[]);

    return this.users;
  }

  public sortBy(testValue: keyof IUser) {
    this.users = this.users.sort((a, b) => a[testValue].localeCompare(b[testValue]));
    return this.users;
  }

  public getTransactions(key: string): ng.IPromise<ITransaction[]> {
    this.transactions = [];
    return this.$http.get<ITransaction[]>(`/rest/user/transactions/${key}`).then((res) => {
      this.transactionsCache[key] = res.data;
      return res.data;
    }, (resp: ng.IHttpResponse<string>) => {

      if (resp.status >= 400 && resp.status < 500) {
        return this.$q.resolve([]);
      }

      return this.$q.reject(resp.data);
    });
  }

  public async displayTransactions(selectedUsers: ReadonlyArray<IUser>) {
    const transactionsPromises = selectedUsers
      .map((row) => row.key)
      .map((key) => {
        const cacheHit = this.transactionsCache[key];
        if (cacheHit) {
          return this.$q.resolve(cacheHit);
        }
        return this.getTransactions(key);
      });

    const transactions = await this.$q.all(transactionsPromises);
    this.transactions = transactions.reduce((acc, val) => acc.concat(val), []); // flatten
  }

  public showMessage(message: string) {
    const toast = this.$mdToast
      .simple()
      .hideDelay(0)
      .textContent(message)
      .highlightAction(true)
      .action('Dismiss');

    this.$mdToast.show(toast);
  }

}

export const adminPageComponent: IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
