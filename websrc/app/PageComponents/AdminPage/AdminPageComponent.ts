import { IComponentOptions, copy } from 'angular';

export const adminRouterState: angular.ui.IState = {
  template: '<admin-page></admin-page>',
  url: '/admin/',
};

export class AdminPageCtrl {

  public usersBackup: any[] = [];
  public users: any[] = [];
  public transaction: any = null;

  constructor(private $http: ng.IHttpService) {
    $http.get<any>('/rest/user/all').then((res) => {
      this.users = res.data.users;
      this.usersBackup = copy(this.users);

      res.data.keys.forEach((key: string, i: number) => {
        this.users[i].key = key;
      });

      this.sortBy('accessId');
    });
  }

  public makeUsersUnique(testValue: string) {
    let prev: any = {};
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

  public sortBy(testValue: string) {
    this.users = this.users.sort((a, b) => a[testValue].localeCompare(b[testValue]));
    return this.users;
  }

  public getTransactions(key: string): ng.IPromise<any> {
    this.transaction = '';
    return this.$http.get(`/rest/user/transactions/${key}`).then(res => {
      this.transaction = res.data;
      return res.data;
    });
  }
}

export const AdminPageComponent: IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
