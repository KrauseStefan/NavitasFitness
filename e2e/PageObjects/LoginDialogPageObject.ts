import { DialogPageObject } from './DialogPageObject';

export class LoginDialogPageObject extends DialogPageObject {

  public loginButton = this.formContainer.$('button[ng-click="submit()"]');
  public resetButton = this.formContainer.$('button[ng-click="openResetPasswordDialog($event)"]');

  public closeButton = this.formContainer.$('.md-toolbar-tools button[ng-click="cancel()"]');
  public cancelButton = this.formContainer.$('md-dialog-actions button[ng-click="cancel()"]');

  public errorEmailNotVerified = this.formContainer.$('[ng-message="emailNotVerified"]');
  public errorCredentialsInvalid = this.formContainer.$('[ng-message="credentialsInvalid"]');

  public openResetForm() {
    this.resetButton.click();
    return new ResetPasswordDialogPageObject();
  }
}

export class ResetPasswordDialogPageObject extends DialogPageObject {
  public resetButton = this.formContainer.$('.md-primary');
}
