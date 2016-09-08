import { IUserDTO, UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;

export class RegistrationFormModel implements IUserDTO {
  public email: string = '';
  public emailRepeat: string = '';
  public password: string = '';
  public passwordRepeat: string = '';
  public navitasId: string = '';

  public toUserDTO(): IUserDTO {
    return {
      email: this.email,
      navitasId: this.navitasId,
      password: this.password,
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

  public submit() {
    this.userService.createUser(this.$scope.model.toUserDTO()).then(() => {
      this.$scope.model = new RegistrationFormModel();
      this.$mdDialog.hide();
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }
}
