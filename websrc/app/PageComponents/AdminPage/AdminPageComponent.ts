import { User } from './Users/User';

export const adminRouterState: angular.ui.IState = {
  template: '<admin-page></admin-page>',
  url: '/admin/',
};

export class AdminPageCtrl implements ng.IComponentController {
  public selectedUserKeys: ReadonlyArray<string> = [];

  public updatingDropbox = false;

  private dropboxUpdateSucceded = this.$mdDialog.alert()
    .clickOutsideToClose(true)
    .title('Dropbox updated')
    .textContent(`Dropbox updated.`)
    .ariaLabel('Dropbox updated')
    .ok('OK');

  private dropboxUpdateFailed = this.$mdDialog.alert()
    .clickOutsideToClose(true)
    .title('Dropbox update failed')
    .textContent(`Dropbox update failed.`)
    .ariaLabel('Dropbox update failed')
    .ok('OK');

  constructor(private $http: ng.IHttpService, private $mdDialog: ng.material.IDialogService) { }

  public userSelectionUpdated(selectedUsers: ReadonlyArray<User>) {
    this.selectedUserKeys = selectedUsers.map((row) => row.key);
  }

  public updateDownboxAccessList() {
    this.updatingDropbox = true;
    this.$http.get('rest/export/csv/')
      .finally(() => this.updatingDropbox = false)
      .then(() => this.$mdDialog.show(this.dropboxUpdateSucceded))
      .catch(() => this.$mdDialog.show(this.dropboxUpdateFailed));
  }

}

export const adminPageComponent: ng.IComponentOptions = {
  controller: AdminPageCtrl,
  templateUrl: '/PageComponents/AdminPage/AdminPage.html',
};
