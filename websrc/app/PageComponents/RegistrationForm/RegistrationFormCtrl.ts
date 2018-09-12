import { IUserDTO, UserService } from '../UserService';
import { RegistrationFormModel } from './RegistrationFormModel';

import IDialogService = ng.material.IDialogService;
import IToastService = ng.material.IToastService;

interface IRegistrationError {
  field?: keyof RegistrationFormModel;
  message: string;
  type: 'invalid' | 'unique_constraint';
}

export class RegistrationForm implements ng.IController {

  constructor(
    private $scope: {
      submit: () => void,
      cancel: () => void,
      model: RegistrationFormModel,
      errorMsg: any,
      RegistrationForm: {[field in keyof RegistrationFormModel]: ng.INgModelController } & ng.IFormController,
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
    this.$scope.RegistrationForm.$pending = true;

    this.userService.createUser(this.toUserDTO(this.$scope.model)).then(() => {
      this.$scope.model = new RegistrationFormModel();
      this.$mdDialog.hide();
      this.displayCheckEmailNotice();
    }, (err: ng.IHttpPromiseCallbackArg<IRegistrationError>) => {
      if (err.data.field && err.data.field.length > 0) {
        const formCtrl = this.getFormFieldCtrl(err.data);
        if (err.data.type) {
          formCtrl.$setValidity(err.data.type, false);
        } else {
          formCtrl.$setValidity('serverValidation', false);
        }
      }
    }).finally(() => this.$scope.RegistrationForm.$pending = false);
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

  private getFieldMap(): { [key: string]: string } {
    const initial: { [key: string]: string } = {};

    return Object.keys(this.$scope.RegistrationForm)
      .filter((i) => i[0] !== '$')
      .reduce((acc, i) => {
        acc[i.toLowerCase()] = i;
        return acc;
      }, initial);
  }

  private getFormFieldCtrl(err: IRegistrationError): ng.INgModelController {
    if (!this.$scope.RegistrationForm[err.field]) {
      const errField = err.field.toLowerCase();
      const fieldMap = this.getFieldMap();
      return this.$scope.RegistrationForm[fieldMap[errField]];
    }

    return this.$scope.RegistrationForm[err.field];
  }

  private displayCheckEmailNotice() {
    return this.$mdDialog.show(
      this.$mdDialog.alert()
        .clickOutsideToClose(true)
        .title('Confirmation e-mail sent')
        .textContent(`Please check your e-mail inbox to compleate registration`)
        .ariaLabel('Confirmation e-mail sent')
        .ok('OK'));
  }
}
