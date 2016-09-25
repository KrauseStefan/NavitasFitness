import { IUserDTO, UserService } from '../UserService';

import IDialogService = angular.material.IDialogService;
import IToastService = angular.material.IToastService;
import IHttpPromiseCallbackArg = angular.IHttpPromiseCallbackArg;
import INgModelController = angular.INgModelController;
import IScope = angular.IScope;

const HttpConflict = 409;

export class RegistrationFormModel implements IUserDTO {
  public email: string = '';
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
    private $scope: {
      submit: () => void,
      cancel: () => void,
      model: RegistrationFormModel,
      errorMsg: any,
      RegistrationForm: {
        email: INgModelController
      }
    } & IScope,
    private userService: UserService,
    private $mdDialog: IDialogService,
    private $mdToast: IToastService) {

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    $scope.model = new RegistrationFormModel();
  }

  public submit() {
    this.userService.createUser(this.$scope.model.toUserDTO()).then(() => {
      this.$scope.model = new RegistrationFormModel();
      this.$mdDialog.hide();
    }, (errorResponse: IHttpPromiseCallbackArg<string>) => {
      if (errorResponse.status === HttpConflict) {
        this.$scope.RegistrationForm.email.$setValidity('emailAvailable', false);
      }
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }
}
