import { $ } from 'protractor';

export class LoginDialogPageObject {
  public formContainer = $('md-dialog');

  public fieldEmail = this.formContainer.$('input[ng-model="model.email"]');
  public fieldPassword = this.formContainer.$('input[ng-model="model.password"]');

  public buttonLogin = this.formContainer.$('button[ng-click="submit()"]');
  public buttonCancel = this.formContainer.$('button[ng-click="cancel()"]');

  public errorLoginSuccessful = this.formContainer.$('[ng-message="loginSuccessful"]');
}
