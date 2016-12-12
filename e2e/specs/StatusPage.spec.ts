import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { StatusPageObject } from '../PageObjects/StatusPageObject';
import { verifyBrowserLog } from '../utility';
import { $, browser } from 'protractor';

describe('StatusPage tests', () => {

  const userInfo = {
    email: 'status-test@domain.com',
    navitasId: '1234509876',
    password: 'Password123',
  };

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

  it('should not be able to click status before being logged in', () => {
    NavigationPageObject.statusPageTab.click()
      .then(() => fail(), () => {/* */ });
  });

  it('[META] login user', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();
  });

  it('should be able to click status when logged in', () => {
    NavigationPageObject.statusPageTab.click();
    expect(browser.getCurrentUrl()).toContain('status');
  });

  it('should be able to process a payment', () => {
    StatusPageObject.waitForPaypalSimBtn();
    StatusPageObject.triggerPaypalPayment();

    NavigationPageObject.statusPageTab.click();

    expect($('tr td:nth-child(3)').getText()).toBe('Completed');
  });

});
