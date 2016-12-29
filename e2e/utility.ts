import * as http from 'http';

import { browser } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

export interface IBrowserLog {
  level: {
    name: string, // SERVERE
    value: number
  };
  message: string;
  timestamp: number;
  type: string;
}

export function verifyBrowserLog(expectedEntries: string[] = []) {
  return (<any>browser).manage().logs().get('browser').then((browserLogs: IBrowserLog[]) => {

    const filteredLog = browserLogs.filter((logEntry) => {
      const index = expectedEntries.findIndex((i) => i === logEntry.message);
      expectedEntries.splice(index, 1);

      return index === undefined;
    });

    if (filteredLog.length > 0) {
      const entries = filteredLog
        .map((entry) => `[${browserLogs[0].type}][${browserLogs[0].level.name}] ${browserLogs[0].message}`);
      throw `Error was thrown during test execution:\n [${entries.join('\n')}`;
    }

    if (expectedEntries.length > 0) {
      throw `[Expected log to contain entry, but it did not: ${expectedEntries.join(', ')}]`;
    }
  });
}

const sessionCoockieKey = 'Session-Key';

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
