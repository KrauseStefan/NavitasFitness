import { IComponentOptions } from 'angular';

export const adminRouterState: angular.ui.IState = {
  template: '<admin-page></admin-page>',
  url: '/admin/',
};

export class AdminPageCtrl {

  public users: any[] = [];
  public transaction: any = null;

  constructor(private $http: ng.IHttpService) {
    $http.get<any>('/rest/user/all').then((res) => {
      this.users = res.data.users;
      res.data.keys.forEach((key: string, i: number) => {
        this.users[i].key = key;
      });
    });
  }

  public getTransactions(key: string) {
    this.transaction = '';
    this.$http.get(`/rest/user/transactions/${key}`).then(res => {
      this.transaction = res.data;
    });
  }
}

export const AdminPageComponent: IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
