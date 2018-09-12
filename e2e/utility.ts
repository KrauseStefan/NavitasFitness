import { browser } from 'protractor';

export interface IBrowserLog {
  level: {
    name: string, // SERVERE
    value: number,
  };
  message: string;
  timestamp: number;
  type: string;
}

export function verifyBrowserLog(expectedEntries: string[] = []): Promise<void> {
  return (<any>browser).manage().logs().get('browser').then((browserLogs: IBrowserLog[]) => {

    const filteredLog = browserLogs.filter((logEntry) => {
      const index = expectedEntries.findIndex((i) => i === logEntry.message);
      expectedEntries.splice(index, 1);

      return index === undefined;
    });

    if (filteredLog.length > 0) {
      const entries = filteredLog
        .map((entry) => `[${browserLogs[0].type}][${browserLogs[0].level.name}] ${browserLogs[0].message}`);
      throw new Error(`Error was thrown during test execution:\n [${entries.join('\n')}`);
    }

    if (expectedEntries.length > 0) {
      throw new Error(`[Expected log to contain entry, but it did not: ${expectedEntries.join(', ')}]`);
    }
  });
}

export function waitForPageToLoad(): Promise<{}> {
  function hasPageLoaded(): Promise<{}> {
    return Promise.resolve(browser.executeAsyncScript(() => {
      const callback = arguments[arguments.length - 1];

      if ((<any>window).waitForAngular) {
        const angular = (<any>window).angular;
        const injector = angular && angular.element(document).injector();

        if (injector) {
          const $browser = injector.get('$browser');
          $browser.notifyWhenNoOutstandingRequests(() => callback(true));
        }
      } else {
        callback(document.readyState);
      }

      callback(false);
    }));
  }

  return Promise.resolve(browser.wait(hasPageLoaded, 10000, 'Page did not load'));
}

export function retryCall<T>(fn: () => Promise<T>, count: number): Promise<T> {
  return new Promise((resolve, reject) => {
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
