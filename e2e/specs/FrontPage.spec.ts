import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { FrontPageObject } from '../PageObjects/FrontPageObject';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { Key, browser } from 'protractor';

const userInfo = {
  name: 'test',
  email: 't-frontpage-admin@domain.com',
  accessId: 'AAMS-asa',
  password: 'Password123',
};

describe('Frontpage tests', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] load page', () => {
    browser.get('/');
  });

  it('[META] create user', () => {
    new DataStoreManipulator().removeUserByEmail(userInfo.email).destroy();

    const regDialog = NavigationPageObject.openRegistrationDialog();

    regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      accessId: userInfo.accessId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.buttonRegister.click();
  });

  it('should not be able to click edit before being logged in', () => {
    expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] login user', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();
    DataStoreManipulator.sendValidationRequest(userInfo.email);

    loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();
  });

  it('should not be able to click edit if logged in with a normal user', () => {
    expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] make user admin', () => {
    new DataStoreManipulator().makeUserAdmin(userInfo.email).destroy();
    browser.refresh();
  });

  it('should be able to enter edit mode as admin', () => {
    const testText = 'This is a test message';
    const textPromise = FrontPageObject.editableArea.getText();

    FrontPageObject.adminEditBtn.click();

    textPromise.then(text => {
      const backspaces = new Array(text.length + 1).join(Key.DELETE);
      FrontPageObject.editableArea.sendKeys(backspaces);
    });

    FrontPageObject.editableArea.sendKeys(testText);

    FrontPageObject.adminSaveBtn.click();

    expect(FrontPageObject.editableArea.getText()).toEqual(testText);
  });

  it('[META] user admin', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogout.click();
  });

});
