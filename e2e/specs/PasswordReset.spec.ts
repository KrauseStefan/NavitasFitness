import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { ResetPasswordDialogPageObject } from '../PageObjects/ResetPasswordDialogPageObject';

import { browser } from 'protractor';
import { stringify } from 'querystring';
import { verifyBrowserLog } from '../utility';

describe('Reset password', () => {

  const newPassword = 'Password2';
  const userInfo = {
    name: 'test',
    email: 'reset-test@domain.com',
    accessId: 'reset-test',
    password: 'Password1',
  };

  afterEach(() => verifyBrowserLog());

  it('[META] create user', async () => {
    await browser.get('/');
    await DataStoreManipulator.removeUserByEmail(userInfo.email);
    const regDialog = await NavigationPageObject.openRegistrationDialog();

    await regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      accessId: userInfo.accessId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });

    await regDialog.buttonRegister.click();
    await AlerDialogPageObject.mainButton.click();
    await expect(regDialog.formContainer.isPresent()).toBe(false);
  });

  it('should be able to request a password reset', async () => {
    const loginDialog = await NavigationPageObject.openLoginDialog();
    const resetDialog = await loginDialog.openResetForm();

    await resetDialog.fillForm({
      email: userInfo.email,
    });

    await resetDialog.resetButton.click();
    await AlerDialogPageObject.mainButton.click();
  });

  it('should be able to reset password', async () => {
    const passwordResetKey = await DataStoreManipulator.getUserId(userInfo.email);
    const passwordResetSecret = await DataStoreManipulator.getUserEntityResetSecretFromEmail(userInfo.email);

    const parms = stringify({ passwordResetKey, passwordResetSecret });

    await browser.get('/main-page/?' + parms);
    const resetPasswordDialog = new ResetPasswordDialogPageObject();
    await resetPasswordDialog.fillForm({
      password: newPassword,
      passwordRepeat: newPassword,
    });

    await resetPasswordDialog.resetButton.click();
  });

  it('should fail login with old password', async () => {
    await DataStoreManipulator.performEmailVerification(userInfo.email);
    const loginDialog = await NavigationPageObject.openLoginDialog();

    await loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });
    await loginDialog.loginButton.click();

    await expect(loginDialog.formContainer.isDisplayed()).toBe(true);

    await expect(loginDialog.errorCredentialsInvalid.isDisplayed()).toBe(true);

    await loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to login with new password', async () => {
    const loginDialog = await NavigationPageObject.openLoginDialog();

    await loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: newPassword,
    });
    await loginDialog.loginButton.click();

    await expect(loginDialog.formContainer.isPresent()).toBe(false);
  });

  it('[META] should be able to logout', async () => {
    await NavigationPageObject.menuButton.click();
    await NavigationPageObject.menuLogout.click();
  });

});
