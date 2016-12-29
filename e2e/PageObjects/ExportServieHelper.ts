import * as http from 'http';

import { browser } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

const sessionCoockieKey = 'Session-Key';

export enum columns {
  "SysID" = 1,
  "DateActivation",
  "SysID2",
  "DateStart",
  "DateEnd",
  "TimeScheme",
  "Comments",
}

function getSessionCookie(): wdp.Promise<string> {
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
