import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { FrontPageObject } from '../PageObjects/FrontPageObject';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';

import { browser, Key } from 'protractor';
import { verifyBrowserLog } from '../utility';

const userInfo = {
  name: 'front-a',
  email: 'front-a@domain.com',
  accessId: 'front-a',
  password: 'Password1',
};

describe('Frontpage tests', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] load page', async () => {
    await browser.get('/');
  });

  it('[META] create user', async () => {
    await DataStoreManipulator.loadUserKinds();
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
  });

  it('should not be able to click edit before being logged in', async () => {
    await expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] login user', async () => {
    const loginDialog = await NavigationPageObject.openLoginDialog();
    await DataStoreManipulator.loadUserKinds();
    await DataStoreManipulator.performEmailVerification(userInfo.email);

    await loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });

    await loginDialog.loginButton.click();
  });

  it('should not be able to click edit if logged in with a normal user', async () => {
    await expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] make user admin', async () => {
    await DataStoreManipulator.loadUserKinds();
    await DataStoreManipulator.makeUserAdmin(userInfo.email);
    await browser.refresh();
  });

  it('should be able to enter edit mode as admin', async () => {
    const testText = 'This is a test message';
    const text = await FrontPageObject.editableArea.getText();

    await FrontPageObject.adminEditBtn.click();

    const backspaces = new Array(text.length + 1).join(Key.DELETE);
    await FrontPageObject.editableArea.sendKeys(backspaces);

    await FrontPageObject.editableArea.sendKeys(testText);

    await FrontPageObject.adminSaveBtn.click();

    await expect(FrontPageObject.editableArea.getText()).toEqual(testText);
  });

  it('[META] user admin', async () => {
    await NavigationPageObject.menuButton.click();
    await NavigationPageObject.menuLogout.click();
  });

});
