import { excel } from '../Typings/exceljs';
import * as Excel from 'exceljs';

import * as http from 'http';

import { browser } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

import { IParsedDate, dateParts } from './StatusPageObject';

declare const Excel: {
  Workbook: excel.IWorkbook;
};

const sessionCoockieKey = 'Session-Key';
export const exportServiceUrl = 'http://localhost:8080/rest/export/xlsx';

export interface IExcelRow {
  SysID: string;
  DateActivation: IParsedDate;
  SysID2: string;
  DateStart: IParsedDate;
  DateEnd: IParsedDate;
  TimeScheme: string;
  Comments: string;
}

function getSessionCookie(): wdp.Promise<string> {
  return browser
    .manage()
    .getCookie(sessionCoockieKey)
    .then((cookie) => {
      return `${sessionCoockieKey}=${cookie.value}`;
    });
}

export function parseUrlToReqOj(url: string): http.RequestOptions {
  const [protocol, host, port, path] = (<Array<string>>url.match(/([A-z]*:)\/\/([A-z]*):(\d*)([\/|\w]*)/)).slice(1);

  return {
    host,
    path,
    port: parseInt(port, 10),
    protocol,
  };
}

function sendRequstWithCookie(url: string, Cookie?: string) {
  const options = parseUrlToReqOj(url);
  options.headers = Cookie ? { Cookie } : {};
  options.method = 'get';

  return new wdp.Promise<http.IncomingMessage>((resolve, reject) => {
    const req = http.request(options, resolve);
    req.end();
  });
}

export function makeRequest(url: string, useSession: boolean = false): wdp.Promise<http.IncomingMessage> {
  if (useSession) {
    return getSessionCookie().then((cookie) => {
      return sendRequstWithCookie(url, cookie);
    });
  } else {
    return sendRequstWithCookie(url);
  }
}

export function sendPayment(custom: string, paymentDate: string): wdp.Promise<string> {
  const dataData = {
    cmd: '_s-xclick',
    custom,
    hosted_button_id: 'KE6ZXLGBP6TRQ',
    payment_date: paymentDate,
  };
  const dataStr = Object.keys(dataData).reduce((previousValue, currentValue) => {
    return `${currentValue}=${dataData[currentValue]}&${previousValue}`;
  }, '').slice(0, -1);

  let options = parseUrlToReqOj('http://localhost:8081/processPayment');
  options.method = 'post';
  options.headers = {
    'Content-Type': 'application/x-www-form-urlencoded',
    'Content-Length': dataStr.length,
  };

  return new wdp.Promise<string>((resolve, reject) => {
    let returnedData = '';
    let successStatus = false;

    const req = http.request(options, (res) => {
      if (res.statusCode && res.statusCode <= 200 && res.statusCode < 300) {
        successStatus = true;
      }
      res.setEncoding('utf8');
      res.on('data', (chunk) => {
        returnedData += chunk;
      });
      res.on('close', (resp) => {
        successStatus ? resolve(returnedData) : reject(res.statusMessage + ': ' + returnedData);
      });
      res.on('end', (resp) => {
        successStatus ? resolve(returnedData) : reject(res.statusMessage + ': ' + returnedData);
      });
      res.on('error', reject);
    });
    req.write(dataStr);
    req.end();
  }).then((data) => {
    return browser.sleep(500).then(() => {
      return data;
    });
  });

}

function parseXlsxDocument(resp: http.IncomingMessage): wdp.Promise<excel.IWorkbook> {
  const workbook = new Excel.Workbook();

  const inputStream = workbook.xlsx.createInputStream();
  resp.pipe(inputStream);
  return new wdp.Promise<excel.IWorkbook>((resolve, reject) => {
    inputStream.on('done', (listener) => {
      resolve(workbook);
    });
  });
}

export function downloadXsltTransactionExport(): wdp.Promise<IExcelRow[]> {
  const respP = makeRequest(exportServiceUrl, true);
  const workbookP = respP.then(parseXlsxDocument);
  return workbookP.then((workbook) => {
    const sheetData = workbook.getWorksheet(1).getSheetValues();
    return sheetData.map(data => <IExcelRow>{
      SysID: data[1],
      DateActivation: dateParts(data[2]),
      SysID2: data[3],
      DateStart: dateParts(data[4]),
      DateEnd: dateParts(data[5]),
      TimeScheme: data[6],
      Comments: data[7],
    });
  });
}
