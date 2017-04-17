import { ResetPasswordFormController } from '../ResetPasswordForm/ResetPasswordForm';
import { IBaseUserDTO, UserService } from '../UserService';
import { element } from 'angular';

import IDialogService = angular.material.IDialogService;
import IHttpPromiseCallbackArg = angular.IHttpPromiseCallbackArg;

const HttpUnauthorized = 401;
const StatusForbidden = 403;

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
    const model: IBaseUserDTO = {
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
    this.userService.createUserSession(this.$scope.model).then(() => {
      this.resetForm();
      this.$mdDialog.hide();
    }, (errorResponse: IHttpPromiseCallbackArg<string>) => {
      if (errorResponse.status === HttpUnauthorized) {
        this.$scope.LoginForm.password.$setValidity('credentialsInvalid', false);
      } else if (errorResponse.status === StatusForbidden) {
        this.$scope.LoginForm.accessId.$setValidity('emailNotVerified', false);
      }
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

}
