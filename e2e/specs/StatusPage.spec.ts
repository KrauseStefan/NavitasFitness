import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';

import { IParsedDate, StatusPageObject, TransactionTableCells} from '../PageObjects/StatusPageObject';

describe('Payments', () => {

  const userInfo = {
    name: 'status',
    email: 'status@domain.com',
    accessId: 'status',
    password: 'Password1',
  };

  afterEach(() => verifyBrowserLog());

  describe('StatusPage', () => {

    it('[META] create user', async () => {
      await browser.get('/');
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
      await expect(regDialog.formContainer.isPresent()).toBe(false, 'Registration dialog was not closed');
    });

    it('should not be able to click status before being logged in', async () => {
      await NavigationPageObject.statusPageTab.click()
        .then(() => fail(), () => {/* */ });
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

    it('should be able to click status when logged in', async () => {
      await NavigationPageObject.statusPageTab.click();
      await expect(browser.getCurrentUrl()).toContain('status');
    });

    it('should report an inactive subscription when no payment is made', async () => {
      await expect(StatusPageObject.getStatusMsgFieldValue()).toEqual('inActive');
      await expect(StatusPageObject.getValidUntilFieldValue()).toEqual('-');
    });

    it('should not be able to process a payment before terms has been accepted', async () => {
      await StatusPageObject.waitForPaypalSimBtn();
      await NavigationPageObject.statusPageTab.click();

      await expect(StatusPageObject.paypalBtn.isEnabled()).toBe(false);
    });

    it('should be able to process a payment', async () => {
      await StatusPageObject.waitForPaypalSimBtn();
      await StatusPageObject.termsAcceptedChkBx.click();
      await StatusPageObject.triggerPaypalPayment();
      await NavigationPageObject.statusPageTab.click();
      await browser.wait(() => {
        return StatusPageObject.getTableCellText(1, TransactionTableCells.Status)
          .then((status) => status === 'Completed', () => false);
      }, 10000, 'Payment could not be compleatly processed');
    });

    it('should show subscription active when show subscription end date when subscribed', async () => {
      function diffMonth(start: IParsedDate, end: IParsedDate): number {
        if (start.year === end.year) {
          return end.month - start.month;
        } else if (start.year + 1 === end.year) {
          return (end.month + 12) - start.month;
        }
        throw 'Date invalid';
      }

      const monthDiff = await StatusPageObject.getPageDates()
        .then(dates => diffMonth(dates.firstTrxDate, dates.validUntil));

      await expect(StatusPageObject.getStatusMsgFieldValue()).toEqual('active');
      await expect(monthDiff).toEqual(6);
    });
  });

});
