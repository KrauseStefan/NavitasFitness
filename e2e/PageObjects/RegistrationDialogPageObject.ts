import { DialogPageObject } from './DialogPageObject';

export class RegistrationDialogPageObject extends DialogPageObject {
  public buttonRegister = this.formContainer.$('button[ng-click="submit()"]');

  public closeButton = this.formContainer.$('.md-toolbar-tools button[ng-click="cancel()"]');
  public cancelButton = this.formContainer.$('md-dialog-actions button[ng-click="cancel()"]');

  public errorEmailConflict = this.formContainer.$('[ng-message="unique_constraint"]');
  public errorAccessIdConflict = this.formContainer.$('[ng-message="unique_constraint"]');
  public errorAccessIdInvalid = this.formContainer.$('[ng-message="invalid"]');

  public errorPasswordDifferent = this.formContainer.$('[ng-message="nfShouldEqual"]');
}
