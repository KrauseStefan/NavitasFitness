import { retryCall } from '../utility';
import * as http from 'http';
import { ElementFinder, ProtractorBrowser, browser as mainBrowser, by } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

let browser: ProtractorBrowser;

const accessIdCol = 5;
const emailCol = 8;

export class DataStoreManipulator {

  public static sendValidationRequest(email: string): wdp.Promise<void> {
    const dataStoreManipulator = new DataStoreManipulator();
    const key = dataStoreManipulator.getUserEntityIdFromEmail(email);
    dataStoreManipulator.destroy();

    return DataStoreManipulator.sendValidationRequestFromKey(key);
  }

  public static sendValidationRequestFromKey(key: wdp.Promise<string>): wdp.Promise<void> {

    return key.then((keyStr) => {
      return new wdp.Promise<void>((resolve, reject) => {
        const url = '/rest/user/verify?code=' + keyStr;
        return http.get({
          port: '9000',
          path: url,
        }, (response: http.ClientResponse) => {
          if (response.headers['location'].includes('Verified=true')) {
            resolve(void {});
          } else {
            reject(void {});
          }
        });
      });

    });

  }

  private deleteBtn: ElementFinder;

  constructor() {
    browser = mainBrowser.forkNewDriverInstance(false, false);
    browser.ignoreSynchronization = true;
    browser.driver.get('http://localhost:8000/datastore?kind=User');

    this.deleteBtn = browser.$('#delete_button');
  }

  public destroy() {
    browser.sleep(200);
    browser.quit();
  }

  public getUserEntityIdFromEmail(email: string): wdp.Promise<string> {
    const queryStr = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${emailCol})').text() === '${email}');

      if(row.length <= 0){
        return;
      }
      return row.find('a')[0]
        .href
        .match(/\\/edit\\/([\\w|\\d|\\-|%]*)?/)[1];
      `;

    return browser.driver.executeScript(queryStr).then((key: string) => {
      if (key) {
        return key;
      }
      throw `Unable to lookup user DB key, email used: ${email}`;
    });
  }

  public removeUserByAccessId(accessId: string): DataStoreManipulator {
    this.selecteItem(accessIdCol, accessId);
    this.deleteSelected();

    return this;
  }

  public removeUserByEmail(email: string): DataStoreManipulator {
    this.selecteItem(emailCol, email);
    this.deleteSelected();

    return this;
  }

  public makeUserAdmin(email: string): DataStoreManipulator {
    this.openItem(emailCol, email);

    const selectAdmin = `document.querySelector('select[name="bool|IsAdmin"]').value = 1;`;
    browser.driver.executeScript(selectAdmin);
    browser.$('input[value="Save Changes"]').click();
    return this;
  }

  private deleteSelected() {
    this.deleteBtn.isPresent()
      .then(isPresent => isPresent ? this.deleteBtn.isEnabled() : wdp.fullyResolved<boolean>(false))
      .then(isEnabled => {
        if (isEnabled) {
          this.deleteBtn.click();
          return retryCall(() => browser.switchTo().alert().accept(), 10);
        } else {
          return wdp.fullyResolved<void>({});
        }
      });
  }

  private openItem(column: number, value: string) {
    const clientSideScript = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      return row.find('a')[0];
    `;
    const itemLink = browser.element(by.js(clientSideScript));
    return itemLink.isPresent().then(isPresent => {
      if (isPresent) {
        return itemLink.click();
      }

      return wdp.fullyResolved<void>({});
    });
  }

  private selecteItem(column: number, value: string) {
    const clientSideScript = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      return row.find('input[type="checkbox"]');
    `;
    const itemChkBox = browser.element(by.js(clientSideScript));
    return itemChkBox.isPresent().then(isPresent => {
      if (isPresent) {
        return itemChkBox.click();
      }

      return wdp.fullyResolved<void>({});
    });
  }
}
