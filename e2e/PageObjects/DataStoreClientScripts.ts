import { ElementFinder, ProtractorBrowser } from 'protractor';

export class DataStoreClientScripts {

  public static getProxy(browser: ProtractorBrowser): DataStoreClientScripts {
    this.proxy = Object.getOwnPropertyNames(DataStoreClientScripts.prototype)
      .reduce((acc, name) => {
        acc[name] = (...args) => {
          // console.log(`clientScripts.${name}.apply(clientScripts, `, args, ')');
          return browser.executeScript(`return clientScripts.${name}.apply(clientScripts, arguments)`, ...args);
        };
        return acc;
      }, <DataStoreClientScripts><any>{});

    return this.proxy;
  }
  private static proxy: DataStoreClientScripts | null = null;

  private columnCache: string[] | null = null;

  public getColumnNumber(columnToMatch: string): number {
    const index = this.getColumnCache().indexOf(columnToMatch.toLowerCase());
    if (index < 0) {
      throw new Error(`Could not lookup column: ${columnToMatch}`);
    }

    return index;
  }

  public getRow(columnToMatch: string, matchValue: string): HTMLTableRowElement {
    const columnIndex = this.getColumnNumber(columnToMatch);
    const tdElements = document.querySelectorAll(`.ae-table.ae-settings-block td:nth-child(${columnIndex + 1})`);
    const tdElement =  Array.prototype.slice.call(tdElements)
      .find((col) => col.innerHTML.trim() === matchValue);

    if (!tdElement) {
      throw new Error('Datarow does not exist');
    }

    return tdElement.parentElement as HTMLTableRowElement;
  }

  public getValue(columnToMatch: string, matchValue: string, columnToGet: string) {
    const row = this.getRow(columnToMatch, matchValue);
    const columnIndex = this.getColumnNumber(columnToGet);
    return this.getFieldText(row.children.item(columnIndex));
  }

  public getRowCheckbox(columnToMatch: string, matchValue: string): ElementFinder {
    const row = this.getRow(columnToMatch, matchValue);
    return row.querySelector('input') as any;
  }

  public getRowIdLink(columnToMatch: string, matchValue: string): ElementFinder {
    const row = this.getRow(columnToMatch, matchValue);
    return row.querySelector('a') as any;
  }

  private getFieldText(fieldElement: Element) {
    function inner() {
      const a = fieldElement.querySelector('a');
      if (a) {
        const matches = a.href.match(/edit\/([a-z,A-Z,0-9,-]*)/);
        if (matches && matches[1]) {
          return matches[1];
        }
        return a.innerHTML;
      }
      return fieldElement.innerHTML;
    }
    return inner().trim();
  }

  private getColumnCache(): string[] {
    if (this.columnCache === null) {
      const columnHeaders = document.querySelectorAll('.ae-table.ae-settings-block th');
      return this.columnCache = Array.prototype.slice.call(columnHeaders)
        .map(this.getFieldText)
        .map((str) => str.toLowerCase());
    }
    return this.columnCache;
  }
}
