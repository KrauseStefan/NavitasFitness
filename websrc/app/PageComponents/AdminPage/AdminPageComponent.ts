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
  private selectedUsers: IUser[] = [];
  private gridApiDefered = this.$q.defer<uiGrid.IGridApiOf<IUser>>();
  private gridApiPromise = this.gridApiDefered.promise;

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
    }, []);

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

  public async displayTransactions(selectedUsers: IUser[]) {
    this.selectedUsers = selectedUsers;
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

  public mergeSelected() {
    const userKeys = this.selectedUsers.map((i) => i.key).join(';');

    this.$http.post(`/rest/user/merge/${userKeys}`, {})
      .then(() => this.$mdToast.show(this.$mdToast.simple().textContent('Success')))
      .catch((resp) => {
        this.showMessage(resp.data);
        return this.$q.resolve();
      });
  }

  public deleteInactiveUsers() {
    this.filterInactiveUsers()
      .then((keys) => {
        if (keys.length > 0) {
          const usersToDelete = keys.join(';');
          return this.$http.delete(`/rest/user/${usersToDelete}`)
            .then((resp) => resp.data);
        } else {
          return 'No duplicated users found';
        }
      })
      .then((data) => this.showMessage(JSON.stringify(data)))
      .catch((data) => {
        this.showMessage(data);
        return this.$q.resolve();
      });
  }

  public deselectActiveUsers() {
    this.$q.all({ keys: this.filterInactiveUsers(), gridApi: this.gridApiPromise })
      .then((data) => {
        const selectionGridApi = data.gridApi.selection;
        const keys = data.keys;
        const usersToSelect = this.selectedUsers
          .filter((user) => keys.indexOf(user.key) !== -1);

        selectionGridApi.clearSelectedRows();

        usersToSelect.forEach((user) => {
          selectionGridApi.selectRow(user);
        });

        this.displayTransactions(usersToSelect);
      })
      .catch((data) => {
        this.showMessage(data);
        return this.$q.resolve();
      });
  }

  private filterInactiveUsers(): ng.IPromise<string[]> {
    const keys = this.selectedUsers.map((i) => i.key).join(';');

    return this.$http.get<string[]>(`/rest/user/duplicated-inactive/${keys}`)
      .then((resp) => resp.data);
  }

}

export const AdminPageComponent: IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
