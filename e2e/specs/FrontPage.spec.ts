import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { FrontPageObject } from '../PageObjects/FrontPageObject';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { Key, browser } from 'protractor';

const userInfo = {
  email: 'frontpage-test-admin@domain.com',
  navitasId: '1234509876',
  password: 'Password123',
};

describe('Navigation tests', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] create user', () => {
    new DataStoreManipulator().removeUser(userInfo.email).destroy();
    browser.get('/');

    const regDialog = NavigationPageObject.openRegistrationDialog();

    regDialog.fillForm({
      email: userInfo.email,
      navitasId: userInfo.navitasId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.termsAcceptedChkBx.click();
    regDialog.buttonRegister.click();
  });

  it('should not be able to click edit before being logged in', () => {
    expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] login user', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();
  });

  it('should not be able to click edit if logged in with a normal user', () => {
    expect(FrontPageObject.adminEditBtn.isPresent()).toBe(false);
  });

  it('[META] user admin', () => {
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
