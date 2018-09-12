import { UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;

export const KeyFieldName = 'passwordResetKey';
export const secretFieldName = 'passwordResetSecret';

export class PasswordChangeFormCtrl {

  private resetKey: string;
  private resetSecret: string;

  constructor(
    private $scope: any,
    private $location: ng.ILocationService,
    private userService: UserService,
    private $mdDialog: IDialogService) {

    const searchParams = this.$location.search();
    this.resetKey = searchParams[KeyFieldName];
    this.resetSecret = searchParams[secretFieldName];

    $scope.model = {
      password: '',
    };

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
  }

  public submit() {
    const dto = {
      password: this.$scope.model.password,
      key: this.resetKey,
      secret: this.resetSecret,
    };

    this.userService.sendPasswordChangeRequest(dto).then(
      () => {
        this.$mdDialog.hide();
      }, (err) => {
        if (err.status >= 400 && err.status < 500) {
          this.showInvalidResetDataAlert();
        } else {
          this.$scope.ResetPasswordForm.password.$setValidity('UnknownError', false);
          this.$scope.ResetPasswordForm.passwordRepeat.$setValidity('UnknownError', false);
        }
      });
  }

  public cancel() {
    this.$location.search({});
    this.$mdDialog.cancel();
  }

  private showInvalidResetDataAlert() {
    const alert = this.$mdDialog.alert()
      .title('Invlaid reset data')
      .textContent('The reset url used is not valid, a new reset request must be made.')
      .ok('Close');

    this.$mdDialog
      .show(alert)
      .finally(() => {
        this.$location.search({});
      });
  }
}
