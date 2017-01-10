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

export function waitForPageToLoad(): wdp.Promise<void> {
  return browser.wait(browser.executeScript(() => document.readyState), 1000, 'Page did not load');
}

export function retryCall<T>(fn: () => wdp.Promise<T>, count: number): wdp.Promise<T> {
  return new wdp.Promise((resolve, reject) => {
    function doCall(actualCount: number) {
      fn().then((value) => resolve(value), (error) => {
        if (count > 0) {
          doCall(actualCount - 1);
        } else {
          reject(error);
        }
      });
    }

    doCall(count - 1);
  });
}
