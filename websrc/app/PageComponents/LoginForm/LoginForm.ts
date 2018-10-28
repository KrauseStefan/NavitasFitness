import { element } from 'angular';
import { ResetPasswordFormController } from '../ResetPasswordForm/ResetPasswordForm';
import { BaseUserDTO, UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;
import IHttpPromiseCallbackArg = angular.IHttpPromiseCallbackArg;

const httpUnauthorized = 401;
const statusForbidden = 403;

export class LoginForm {

  constructor(
    private $scope: any,
    private userService: UserService,
    private $mdDialog: IDialogService) {

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    $scope.openResetPasswordDialog = (event: MouseEvent) => this.openResetPasswordDialog(event);
    this.resetForm();
  }

  public resetForm() {
    const model: BaseUserDTO = {
      accessId: '',
      password: '',
    };
    this.$scope.model = model;
  }

  public openResetPasswordDialog(event: MouseEvent) {
    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: ResetPasswordFormController,
      fullscreen: false,
      // multiple: true,
      parent: element(document.body),
      targetEvent: event,
      templateUrl: '/PageComponents/ResetPasswordForm/ResetPasswordForm.html',
    });
  }

  public submit() {
    this.userService.doUserLogin(this.$scope.model).then(() => {
      this.resetForm();
      this.$mdDialog.hide();
    }, (errorResponse: IHttpPromiseCallbackArg<string>) => {
      if (errorResponse.status === httpUnauthorized) {
        this.$scope.LoginForm.password.$setValidity('credentialsInvalid', false);
      } else if (errorResponse.status === statusForbidden) {
        this.$scope.LoginForm.accessId.$setValidity('emailNotVerified', false);
      }
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

}
