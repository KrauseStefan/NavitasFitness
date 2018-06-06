import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

import {
  IExcelRow, downloadXsltTransactionExport, exportServiceUrl, makeRequest, sendPayment,
} from '../PageObjects/ExportServieHelper';
import {
  IParsedDate,
  StatusPageObject as pageObject,
  TransactionTableCells,
} from '../PageObjects/StatusPageObject';

async function getPageDatesAsExportedRow(id: string, email: string): wdp.Promise<IExcelRow> {
  const dates = await pageObject.getPageDates();

  return <IExcelRow>{
    SysID: id,
    DateActivation: dates.firstTrxDate,
    SysID2: id,
    DateStart: dates.firstTrxDate,
    DateEnd: dates.validUntil,
    TimeScheme: "24 Timers",
    Comments: email,
  };
}

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
      await DataStoreManipulator.sendValidationRequest(userInfo.email);

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
      await expect(pageObject.getStatusMsgFieldValue()).toEqual('inActive');
      await expect(pageObject.getValidUntilFieldValue()).toEqual('-');
    });

    it('should not be able to process a payment before terms has been accepted', async () => {
      await pageObject.waitForPaypalSimBtn();
      await NavigationPageObject.statusPageTab.click();

      await expect(pageObject.paypalBtn.isEnabled()).toBe(false);
    });

    it('should be able to process a payment', async () => {
      await pageObject.waitForPaypalSimBtn();
      await pageObject.termsAcceptedChkBx.click();
      await pageObject.triggerPaypalPayment();
      await NavigationPageObject.statusPageTab.click();
      await browser.wait(() => {
        return pageObject.getTableCellText(1, TransactionTableCells.Status)
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

      const monthDiff = await pageObject.getPageDates()
        .then(dates => diffMonth(dates.firstTrxDate, dates.validUntil));

      await expect(pageObject.getStatusMsgFieldValue()).toEqual('active');
      await expect(monthDiff).toEqual(6);
    });
  });

  describe('xlsx export', () => {

    it('should return 401 if not logged to export data', async () => {
      const includeLoginSession = false;
      const resp = await makeRequest(exportServiceUrl, includeLoginSession);

      await expect(resp.statusCode).toBe(401);
    });

    it('should return 401 if user does not have admin rights', async () => {
      const includeLoginSession = true;
      const resp = await makeRequest(exportServiceUrl, includeLoginSession);

      await expect(resp.statusCode).toBe(401);
    });

    it('should be possible to download an xlsx with active subscriptions', async () => {
      await DataStoreManipulator.loadUserKinds();
      await DataStoreManipulator.makeUserAdmin(userInfo.email);

      const pageDates = await getPageDatesAsExportedRow(userInfo.accessId, userInfo.email);
      const userRows = (await downloadXsltTransactionExport())
        .filter(row => row.Comments === userInfo.email);
      const userRow = userRows[0];

      await expect(userRows.length).toBe(1);
      await expect(pageDates).toEqual(userRow);
    });

    it('test', async () => {
      // payment_date: '00:40:46 Jan 01, 2018 CET',
      await sendPayment(userInfo.email, '15:40:46 Dec 31, 2017 PST');

      await NavigationPageObject.mainPageTab.click();
      await NavigationPageObject.statusPageTab.click();

      const pageDates = await getPageDatesAsExportedRow(userInfo.accessId, userInfo.email);
      const userRows = (await downloadXsltTransactionExport())
        .filter(row => row.Comments === userInfo.email);
      const userRow = userRows[0];

      expect(userRows.length).toBe(1);
      expect(pageDates).toEqual(userRow);
    });
  });
});
