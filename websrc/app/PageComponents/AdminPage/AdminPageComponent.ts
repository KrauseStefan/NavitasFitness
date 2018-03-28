import { IComponentOptions, copy } from 'angular';

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
  public transaction: ITransaction[] = [];

  public gridOptions: uiGrid.IGridOptionsOf<IUser> = {
    data: [],
    enableColumnMenus: false,
    enableFiltering: true,
    enableHorizontalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    enableVerticalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    enableRowHeaderSelection: false,
    enableRowSelection: true,
    modifierKeysToMultiSelect: false,
    multiSelect: false,
    noUnselect: true,
    onRegisterApi: (gridApi) => {
      gridApi.selection.on.rowSelectionChanged(
        this.$scope,
        row => this.getTransactions(row.entity.key)
      );
    },
    rowHeight: 42,
  };

  constructor(
    private $scope: ng.IScope,
    private $http: ng.IHttpService,
    private uiGridConstants: uiGrid.IUiGridConstants
  ) {
    const headerHeight = 33;
    const filterHeight = 28;
    $http.get<{ users: IUser[], keys: string[] }>('/rest/user/all').then((res) => {
      this.users = res.data.users;
      this.usersBackup = copy(this.users);
      this.gridOptions.data = this.users;

      // small hack, I suspect this should have been public
      (<any>this.gridOptions).headerRowHeight = Math.ceil((filterHeight + headerHeight) / 2);

      this.gridOptions.minRowsToShow = 10;
      if (this.users.length < this.gridOptions.minRowsToShow) {
        this.gridOptions.minRowsToShow = this.users.length;
        this.gridOptions.enableVerticalScrollbar = this.uiGridConstants.scrollbars.NEVER;
      }

      res.data.keys.forEach((key: string, i: number) => {
        this.users[i].key = key;
      });

      this.sortBy('accessId');
    });
  }

  public makeUsersUnique(testValue: keyof IUser) {
    let prev: IUser;
    this.users = this.users.reduce((acc, value) => {
      if (prev[testValue] === value[testValue]) {
        if (acc.length > 0 && acc[acc.length - 1] !== prev) {
          acc.push(prev);
        }

        acc.push(value);
      }

      prev = value;
      return acc;
    }, []);

    return this.users;
  }

  public sortBy(testValue: keyof IUser) {
    this.users = this.users.sort((a, b) => a[testValue].localeCompare(b[testValue]));
    return this.users;
  }

  public getTransactions(key: string): ng.IPromise<ITransaction[]> {
    this.transaction = [];
    return this.$http.get<ITransaction[]>(`/rest/user/transactions/${key}`).then(res => {
      this.transaction = res.data;
      return res.data;
    });
  }
}

export const AdminPageComponent: IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
