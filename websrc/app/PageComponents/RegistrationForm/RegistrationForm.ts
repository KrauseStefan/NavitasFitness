import { UserService, UserDTO } from '../UserService';

import IDialogService = angular.material.IDialogService;

export class RegistrationFormModel implements UserDTO {
  email: string = "";
  emailRepeat: string = "";
  password: string = "";
  passwordRepeat: string = "";
  navitasId: string = "";

  toUserDTO(): UserDTO {
    return {
      email: this.email,
      password: this.password,
      navitasId: this.navitasId
    };
  }
}

export class RegistrationForm {

  constructor(
    private $scope: any,
    private userService: UserService,
    private $mdDialog: IDialogService) {
    
    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    $scope.model = new RegistrationFormModel();

  }

  submit() {
    this.userService.createUser(this.$scope.model.toUserDTO()).then(() => {
      this.$scope.model = new RegistrationFormModel();
      this.$mdDialog.hide();
    });
  }

  cancel() {
    this.$mdDialog.cancel();
  }
}