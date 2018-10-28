import { User } from './Users/User';

export const adminRouterState: angular.ui.IState = {
  template: '<admin-page></admin-page>',
  url: '/admin/',
};

export class AdminPageCtrl implements ng.IComponentController {
  public selectedUserKeys: ReadonlyArray<string> = [];

  public userSelectionUpdated(selectedUsers: ReadonlyArray<User>) {
    this.selectedUserKeys = selectedUsers.map((row) => row.key);
  }
}

export const adminPageComponent: ng.IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
