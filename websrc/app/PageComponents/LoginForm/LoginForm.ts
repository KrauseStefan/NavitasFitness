import { UserService, BaseUserDTO } from '../UserService';

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
