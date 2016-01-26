/// <reference path=".../../../../../typings/angularjs/angular.d.ts"/>
/// <reference path=".../../../../../typings/angular-material/angular-material"/>

import { UserService, BaseUserDTO } from '../UserService'

export class LoginForm {
  
  constructor(
    private $scope: any,
    private userService: UserService,
    private $mdDialog: angular.material.IDialogService) {
        
    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    this.resetForm();
  }
  
  resetForm() {
    var model: BaseUserDTO = {
      email: '',
      password: ''
    };
    this.$scope.model = model;
  }

  submit() {
    this.userService.createUserSession(this.$scope.model).then(() => {
      this.resetForm();
      this.$mdDialog.hide();
    });
  }

  cancel() {
    this.$mdDialog.cancel();
  }

}
