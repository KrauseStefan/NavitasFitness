import { excel } from '../Typings/exceljs';

import * as Excel from 'exceljs';
import * as http from 'http';

import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { StatusPageObject as pageObject } from '../PageObjects/StatusPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';
import { promise as wdpromise } from 'selenium-webdriver';

declare const Excel: {
  Workbook: excel.IWorkbook;
};

function dataParts(date: string): { day: number, month: number, year: number } {
  // format 'DD-MM-YYYY'
  const [day, month, year] = date
    .split('-')
    .map((i) => parseInt(i, 10));

  return {
    day,
    month,
    year,
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

const sessionCoockieKey = 'Session-Key';

function getSessionCookie(): wdpromise.Promise<string> {
  return browser
    .manage()
    .getCookie(sessionCoockieKey)
    .then((cookie) => {
      return `${sessionCoockieKey}=${cookie.value}`;
    });
}

function sendRequstWithCookie(url: string, Cookie?: string) {
  const [protocol, host, port, path] = (<Array<string>>url.match(/([A-z]*:)\/\/([A-z]*):(\d*)([\/|\w]*)/)).slice(1);

  const options: http.RequestOptions = {
    headers: Cookie ? { Cookie } : {},
    host,
    method: 'get',
    path,
    port: parseInt(port, 10),
    protocol,
  };

  return new wdpromise.Promise<http.IncomingMessage>((resolve, reject) => {
    const req = http.request(options, resolve);
    req.end();
  });
}

function makeRequest(url: string, useSession: boolean = false): wdpromise.Promise<http.IncomingMessage> {
  if (useSession) {
    return getSessionCookie().then((cookie) => {
      return sendRequstWithCookie(url, cookie);
    });
  } else {
    return sendRequstWithCookie(url);
  }
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

  it('should return 401 if not logged to export data', () => {
    const statusCodeP = makeRequest('http://localhost:8080/rest/export/xlsx', false).then((resp) => {
      return resp.statusCode;
    });

    expect(statusCodeP).toBe(401);
  });

  it('should return 401 if user does not have admin rights', () => {
    const statusCodeP = makeRequest('http://localhost:8080/rest/export/xlsx', true).then((resp) => {
      return resp.statusCode;
    });

    expect(statusCodeP).toBe(401);
  });

  it('should be possible to download an xslt with active subscriptions', () => {

    new DataStoreManipulator().makeUserAdmin(userInfo.email).destroy();

    const respP = makeRequest('http://localhost:8080/rest/export/xlsx', true);

    const workbookP = respP.then((resp) => {
      const workbook = new Excel.Workbook();

      const inputStream = workbook.xlsx.createInputStream();
      resp.pipe(inputStream);
      return new wdpromise.Promise<excel.IWorkbook>((resolve, reject) => {
        inputStream.on('done', (listener) => {
          resolve(workbook);
        });
      });
    });

    const statusCodeP = respP.then((resp) => {
      return resp.statusCode;
    });

    workbookP.then((workbook) => {

      enum columns {
        "SysID" = 1,
        "DateActivation",
        "SysID2",
        "DateStart",
        "DateEnd",
        "TimeScheme",
        "Comments",
      }

      const worksheet = workbook.getWorksheet(1);
      const userRows = worksheet.getSheetValues().filter(i => i[columns.Comments] === userInfo.email);
      expect(userRows.length).toBe(1);

      const userRow = userRows[0];
      expect(userRow[columns.SysID]).toBe(userInfo.navitasId);
      expect(userRow[columns.SysID2]).toBe(userInfo.navitasId);
      expect(userRow[columns.Comments]).toBe(userInfo.email);

    });

    expect(statusCodeP).toBe(200);
  });
});
