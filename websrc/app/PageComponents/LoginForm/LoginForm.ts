import { IBaseUserDTO, UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;
import IHttpPromiseCallbackArg = angular.IHttpPromiseCallbackArg;

const HttpUnauthorized = 401;

export class LoginForm {

  constructor(
    private $scope: any,
    private userService: UserService,
    private $mdDialog: IDialogService) {

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    this.resetForm();
  }

  public resetForm() {
    const model: IBaseUserDTO = {
      accessId: '',
      password: '',
    };
    this.$scope.model = model;
  }

  public submit() {
    this.userService.createUserSession(this.$scope.model).then(() => {
      this.resetForm();
      this.$mdDialog.hide();
    }, (errorResponse: IHttpPromiseCallbackArg<string>) => {
      if (errorResponse.status === HttpUnauthorized) {
        this.$scope.LoginForm.password.$setValidity('loginSuccessful', false);
      }

    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

}
