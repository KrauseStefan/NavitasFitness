import { DialogPageObject } from './DialogPageObject';

export class LoginDialogPageObject extends DialogPageObject {

  public loginButton = this.formContainer.$('button[ng-click="submit()"]');

  public closeButton = this.formContainer.$('.md-toolbar-tools button[ng-click="cancel()"]');
  public cancelButton = this.formContainer.$('md-dialog-actions button[ng-click="cancel()"]');

  public errorLoginSuccessful = this.formContainer.$('[ng-message="loginSuccessful"]');

}
