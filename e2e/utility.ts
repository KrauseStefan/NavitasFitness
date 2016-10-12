import { browser } from 'protractor';

export interface BrowserLog {
  level: {
    name: string, //SERVERE
    value: number
  },
  message: string,
  timestamp: number,
  type: string
}

export function verifyBrowserLog(expectedEntries: string[] = []) {
  return (<any>browser).manage().logs().get('browser').then((browserLogs: BrowserLog[]) => {

    const filteredLog = browserLogs.filter((logEntry) => {
      const index = expectedEntries.findIndex((i) => i === logEntry.message)
      expectedEntries.splice(index, 1);

      return index === undefined;
    })

    if (filteredLog.length > 0) {
      const entries =filteredLog.map((entry) => `[${browserLogs[0].type}][${browserLogs[0].level.name}] ${browserLogs[0].message}`);
      throw `Error was thrown during test execution:\n [${entries.join('\n')}`;
    }

    if(expectedEntries.length > 0) {
      throw `[Expected log to contain entry, but it did not: ${expectedEntries.join(', ')}]`;
    }
  });
}