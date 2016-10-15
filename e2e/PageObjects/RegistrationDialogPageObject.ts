import { $ } from 'protractor';

export class RegistrationDialogPageObject {
  public formContainer = $('md-dialog');

  public fieldEmail = this.formContainer.$('input[ng-model="model.email"]');
  public fieldPassword = this.formContainer.$('input[ng-model="model.password"]');
  public fieldPasswordRepeat = this.formContainer.$('input[ng-model="model.passwordRepeat"]');
  public fieldNavitasId = this.formContainer.$('input[ng-model="model.navitasId"]');

  public buttonRegister = this.formContainer.$('button[ng-click="submit()"]');
  public buttonCancel = this.formContainer.$('button[ng-click="cancel()"]');

  public errorEmailUnavailable = this.formContainer.$('[ng-message="emailAvailable"]');
}