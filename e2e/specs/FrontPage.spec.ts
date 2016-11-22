import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { FrontPageObject } from '../PageObjects/FrontPageObject';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { Key, browser } from 'protractor';

const userInfo = {
  email: 'test-admin@domain.com',
  navitasId: '1234509876',
  password: 'Password123',
};

fdescribe('Navigation tests', () => {

  afterEach(() => verifyBrowserLog());

  it('should create an admin user', () => {
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
    new DataStoreManipulator().makeUserAdmin(userInfo.email).destroy();

    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();
  });

  it('should be able to enter edit mode', () => {
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

});
