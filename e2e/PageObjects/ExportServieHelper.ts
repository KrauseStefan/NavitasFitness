import { Workbook } from 'exceljs';

import * as http from 'http';

import { browser } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

import { IParsedDate, dateParts } from './StatusPageObject';

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

async function getSessionCookie(): wdp.Promise<string> {
  const cookie = await browser
    .manage()
    .getCookie(sessionCoockieKey);

  return `${sessionCoockieKey}=${cookie.value}`;
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

function sendRequstWithCookie(url: string, Cookie?: string): Promise<http.IncomingMessage> {
  const options = parseUrlToReqOj(url);
  options.headers = Cookie ? { Cookie } : {};
  options.method = 'get';

  return new Promise<http.IncomingMessage>((resolve, reject) => {
    const req = http.request(options, resolve);
    req.end();
  });
}

export async function makeRequest(url: string, useSession: boolean = false): Promise<http.IncomingMessage> {
  if (useSession) {
    const cookie = await getSessionCookie();
    return await sendRequstWithCookie(url, cookie);
  } else {
    return await sendRequstWithCookie(url);
  }
}

export async function sendPayment(custom: string, paymentDate: string): Promise<string> {
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

  const data = await new Promise<string>((resolve, reject) => {
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
  });

  return data;
}

function parseXlsxDocument(resp: http.IncomingMessage): wdp.Promise<Workbook> {
  const workbook = new Workbook();

  const inputStream = workbook.xlsx.createInputStream();
  resp.pipe(inputStream);
  return new Promise<Workbook>((resolve, reject) => {
    inputStream.on('done', (listener) => {
      resolve(workbook);
    });
  });
}

export async function downloadXsltTransactionExport(): Promise<IExcelRow[]> {
  const resp = await makeRequest(exportServiceUrl, true);

  const workbook = await parseXlsxDocument(resp);

  const sheetData = (<any>(workbook.getWorksheet(1))).getSheetValues();
  return sheetData.map(data => <IExcelRow>{
    SysID: data[1],
    DateActivation: dateParts(data[2]),
    SysID2: data[3],
    DateStart: dateParts(data[4]),
    DateEnd: dateParts(data[5]),
    TimeScheme: data[6],
    Comments: data[7],
  });
}
