import { IBaseUserDTO, UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;

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
      email: '',
      password: '',
    };
    this.$scope.model = model;
  }

  public submit() {
    this.userService.createUserSession(this.$scope.model).then(() => {
      this.resetForm();
      this.$mdDialog.hide();
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

}
