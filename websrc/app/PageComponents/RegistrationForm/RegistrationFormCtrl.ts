import { IUserDTO, UserService } from '../UserService';

import IDialogService = ng.material.IDialogService;
import IToastService = ng.material.IToastService;

interface IRegistrationError {
  field?: keyof RegistrationFormModel;
  message: string;
  type: 'invalid' | 'unique_constraint';
}

export class RegistrationFormModel implements IUserDTO {
  public name = '';
  public email = '';
  public password = '';
  public passwordRepeat = '';
  public accessId = '';
}

export class RegistrationForm {

  constructor(
    private $scope: {
      submit: () => void,
      cancel: () => void,
      model: RegistrationFormModel,
      errorMsg: any,
      RegistrationForm: {[field in keyof RegistrationFormModel]: ng.INgModelController }
    } & ng.IScope,
    private userService: UserService,
    private $mdDialog: IDialogService,
    private $mdToast: IToastService) {

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
    $scope.model = new RegistrationFormModel();
  }

  public toUserDTO(formModel: RegistrationFormModel): IUserDTO {
    return {
      name: formModel.name,
      email: formModel.email,
      accessId: formModel.accessId,
      password: formModel.password,
    };
  }

  public submit() {
    this.userService.createUser(this.toUserDTO(this.$scope.model)).then(() => {
      this.$scope.model = new RegistrationFormModel();
      this.$mdDialog.hide();
      this.displayCheckEmailNotice();
    }, (err: ng.IHttpPromiseCallbackArg<IRegistrationError>) => {
      if (err.data.field && err.data.field.length > 0) {
        if (err.data.type) {
          this.$scope.RegistrationForm[err.data.field].$setValidity(err.data.type, false);
        } else {
          this.$scope.RegistrationForm[err.data.field].$setValidity('serverValidation', false);
        }
      }
    });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

  private displayCheckEmailNotice() {
    return this.$mdDialog.show(
      this.$mdDialog.alert()
        .clickOutsideToClose(true)
        .title('Confirmation e-mail sent')
        .textContent(`Please check your e-mail inbox to compleate registration`)
        .ariaLabel('Confirmation e-mail sent')
        .ok('OK')
    );
  }
}
