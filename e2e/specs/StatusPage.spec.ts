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

export function getPageDatesAsExportedRow(id: string, email: string): wdp.Promise<IExcelRow> {
  return pageObject.getPageDates().then((dates) => {
    return <IExcelRow>{
      SysID: id,
      DateActivation: dates.firstTrxDate,
      SysID2: id,
      DateStart: dates.firstTrxDate,
      DateEnd: dates.validUntil,
      TimeScheme: "24 Timers",
      Comments: email,
    };
  });
}

describe('Payments', () => {

  const userInfo = {
    name: 'test',
    email: 'status-test@domain.com',
    accessId: '1234509876',
    password: 'Password1',
  };

  afterEach(() => verifyBrowserLog());

  describe('StatusPage', () => {

    it('[META] create user', () => {
      new DataStoreManipulator().removeUser(userInfo.email).destroy();
      browser.get('/');

      const regDialog = NavigationPageObject.openRegistrationDialog();

      regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: userInfo.accessId,
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
      expect(pageObject.getStatusMsgFieldValue()).toEqual('inActive');
      expect(pageObject.getValidUntilFieldValue()).toEqual('-');
    });

    it('should be able to process a payment', () => {
      pageObject.waitForPaypalSimBtn();
      pageObject.triggerPaypalPayment();
      NavigationPageObject.statusPageTab.click();

      expect(pageObject.getTableCellText(1, TransactionTableCells.Status)).toBe('Completed');
    });

    it('should show subscription active when show subscription end date when subscribed', () => {
      function diffMonth(start: IParsedDate, end: IParsedDate): number {
        if (start.year === end.year) {
          return end.month - start.month;
        } else if (start.year + 1 === end.year) {
          return (end.month + 12) - start.month;
        }
        throw 'Date invalid';
      }

      const monthDiffP = pageObject.getPageDates()
        .then(dates => diffMonth(dates.firstTrxDate, dates.validUntil));

      expect(pageObject.getStatusMsgFieldValue()).toEqual('active');
      expect(monthDiffP).toEqual(6);
    });
  });

  describe('xlsx export', () => {

    it('should return 401 if not logged to export data', () => {
      const includeLoginSession = false;
      const statusCodeP = makeRequest(exportServiceUrl, includeLoginSession)
        .then((resp) => resp.statusCode);

      expect(statusCodeP).toBe(401);
    });

    it('should return 401 if user does not have admin rights', () => {
      const includeLoginSession = true;
      const statusCodeP = makeRequest(exportServiceUrl, includeLoginSession)
        .then((resp) => resp.statusCode);

      expect(statusCodeP).toBe(401);
    });

    it('should be possible to download an xlsx with active subscriptions', () => {
      new DataStoreManipulator().makeUserAdmin(userInfo.email).destroy();

      const pageDatesP = getPageDatesAsExportedRow(userInfo.accessId, userInfo.email);
      const userRowsP = downloadXsltTransactionExport()
        .then(rows => rows.filter(row => row.Comments === userInfo.email));
      const userRowP = userRowsP.then(u => u[0]);

      expect(userRowsP.then(u => u.length)).toBe(1);
      expect(pageDatesP).toEqual(userRowP);
    });

    it('test', () => {
      // payment_date: '00:40:46 Jan 01, 2018 CET',
      sendPayment(userInfo.email, '15:40:46 Dec 31, 2017 PST');

      NavigationPageObject.mainPageTab.click();
      NavigationPageObject.statusPageTab.click();

      const pageDatesP = getPageDatesAsExportedRow(userInfo.accessId, userInfo.email);
      const userRowsP = downloadXsltTransactionExport()
        .then(rows => rows.filter(row => row.Comments === userInfo.email));
      const userRowP = userRowsP.then(u => u[0]);

      expect(userRowsP.then(u => u.length)).toBe(1);
      expect(pageDatesP).toEqual(userRowP);
    });
  });
});
