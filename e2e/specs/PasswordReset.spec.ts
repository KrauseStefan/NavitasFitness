import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { ResetPasswordDialogPageObject } from '../PageObjects/LoginDialogPageObject';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';

import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';
import { stringify } from 'querystring';
import { promise as wdp } from 'selenium-webdriver';

describe('Reset password', () => {

  const newPassword = 'Password2';
  const userInfo = {
    name: 'test',
    email: 'reset-test@domain.com',
    accessId: 'reset-test',
    password: 'Password1',
  };

  afterEach(() => verifyBrowserLog());

  it('[META] create user', () => {
    new DataStoreManipulator().removeUserByEmail(userInfo.email).destroy();
    browser.get('/');

    const regDialog = NavigationPageObject.openRegistrationDialog();

    regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      accessId: userInfo.accessId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.buttonRegister.click();
    AlerDialogPageObject.mainButton.click();
  });

  it('should be request a password rest', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();
    const resetDialog = loginDialog.openResetForm();

    resetDialog.fillForm({
      email: userInfo.email,
    });

    resetDialog.resetButton.click();
    AlerDialogPageObject.mainButton.click();
  });

  it('should be able to reset password', () => {
    const ds = new DataStoreManipulator();

    const passwordResetKey = ds.getUserEntityIdFromEmail(userInfo.email);
    const passwordResetSecret = ds.getUserEntityResetSecretFromEmail(userInfo.email);

    wdp.all([passwordResetKey, passwordResetSecret]).then((array) => {
      const parms = stringify({
        passwordResetKey: array[0],
        passwordResetSecret: array[1],
      });
      ds.destroy();
      browser.get('/main-page/?' + parms);
      const resetPasswordDialog = new ResetPasswordDialogPageObject();
      resetPasswordDialog.fillForm({
        password: newPassword,
        passwordRepeat: newPassword,
      });

      resetPasswordDialog.resetButton.click();
    });
  });

  it('should fail login with old password', () => {
    DataStoreManipulator.sendValidationRequest(userInfo.email);
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });
    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    expect(loginDialog.errorCredentialsInvalid.isDisplayed()).toBe(true);

    loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to login with new password', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: newPassword,
    });
    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isPresent()).toBe(false);
  });

});
