import { DialogPageObject } from './DialogPageObject'

export class RegistrationDialogPageObject extends DialogPageObject {
  public buttonRegister = this.formContainer.$('button[ng-click="submit()"]');

  public closeButton = this.formContainer.$('.md-toolbar-tools button[ng-click="cancel()"]');
  public cancelButton = this.formContainer.$('md-dialog-actions button[ng-click="cancel()"]');

  public errorEmailUnavailable = this.formContainer.$('[ng-message="emailAvailable"]');
}