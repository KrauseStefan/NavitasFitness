import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { StatusPageObject as pageObject } from '../PageObjects/StatusPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';
import { promise as wdpromise } from 'selenium-webdriver';

function dataParts(date: string): { day: number, month: number, year: number } {
  // format 'DD-MM-YYYY'
  const parts = date
    .split('-')
    .map((i) => parseInt(i, 10));

  return {
    day: parts[0],
    month: parts[1],
    year: parts[2],
  };
}

function diffMonth(data) {
  const [start, end] = data;
  if (start.year === end.year) {
    return end.month - start.month;
  } else if (start.year + 1 === end.year) {
    return (end.month + 12) - start.month;
  }
  throw 'Date invalid';
}

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

  it('should report an inactive subscription when no payment is made', () => {
    expect(pageObject.statusMsgField.evaluate('$ctrl.model.statusMsgKey')).toEqual('inActive');
    expect(pageObject.subscriptionEndField.evaluate('$ctrl.model.validUntill')).toEqual('-');
  });

  it('should be able to process a payment', () => {
    pageObject.waitForPaypalSimBtn();
    pageObject.triggerPaypalPayment();

    NavigationPageObject.statusPageTab.click();

    expect(pageObject.getFirstRowCell(3).getText()).toBe('Completed');
  });

  it('should show subscription active when show subscription end date when subscribed', () => {
    const startP = pageObject.getFirstRowCell(2).getText().then(dataParts);
    const endP = pageObject.subscriptionEndField.evaluate('$ctrl.model.validUntill').then(dataParts);
    const monthDiffP = wdpromise.all([startP, endP]).then(diffMonth);

    expect(pageObject.statusMsgField.evaluate('$ctrl.model.statusMsgKey')).toEqual('active');
    expect(monthDiffP).toEqual(6);
  });

});
