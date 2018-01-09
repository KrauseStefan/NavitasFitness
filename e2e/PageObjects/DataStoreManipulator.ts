import { retryCall, waitForPageToLoad } from '../utility';
import * as http from 'http';
import { ElementFinder, ProtractorBrowser, browser as mainBrowser, by } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

let browser: ProtractorBrowser;

const accessIdCol = 5;
const emailCol = 7;
const resetSecretCol = 11;

export class DataStoreManipulator {

  public static async sendValidationRequest(email: string): Promise<void> {
    await DataStoreManipulator.init();
    const key = await DataStoreManipulator.getUserEntityIdFromEmail(email);
    await DataStoreManipulator.destroy();

    await DataStoreManipulator.sendValidationRequestFromKey(key);
  }

  public static sendValidationRequestFromKey(key: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const url = '/rest/user/verify?code=' + key;
      return http.get({
        port: '8080',
        path: url,
      }, (response: http.ClientResponse) => {
        const location = response.headers.location;
        if (location && location.includes('Verified=true')) {
          resolve();
          return;
        }

        reject();
      });
    });
  }

  public static async init(): Promise<void> {
    browser = await mainBrowser.forkNewDriverInstance(false, false, false);
    await browser.waitForAngularEnabled(false);
    await browser.get('http://localhost:8000/datastore?kind=User');
    await browser.ready;

    await browser.waitForAngularEnabled(false);

    this.deleteBtn = await browser.$('#delete_button');
  }

  public static async destroy() {
    await waitForPageToLoad();
    await browser.close();
    await browser.quit();
  }

  public static getUserEntityIdFromEmail(email: string): wdp.Promise<string> {
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

  public static getUserEntityResetSecretFromEmail(email: string): wdp.Promise<string> {
    const queryStr = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${emailCol})').text() === '${email}');

      if(row.length <= 0){
        return;
      }
      return row.find('td')[${resetSecretCol}].innerText;
      `;

    return browser.driver.executeScript(queryStr).then((secret: string) => {
      if (secret) {
        return secret;
      }

      throw `Unable to lookup reset secret, email used: ${email}`;
    });
  }

  public static async removeUserByAccessId(accessId: string): Promise<void> {
    const elementSelected = await DataStoreManipulator.selecteItem(accessIdCol, accessId);
    if (elementSelected) {
      await DataStoreManipulator.deleteSelected();
    }
  }

  public static async removeUserByEmail(email: string): Promise<void> {
    const elementSelected = await DataStoreManipulator.selecteItem(emailCol, email);

    if (elementSelected) {
      await DataStoreManipulator.deleteSelected();
    }
  }

  public static async makeUserAdmin(email: string): Promise<void> {
    await DataStoreManipulator.openItem(emailCol, email);

    const selectAdmin = `document.querySelector('select[name="bool|IsAdmin"]').value = 1;`;
    await browser.driver.executeScript(selectAdmin);
    await browser.$('input[value="Save Changes"]').click();
  }

  private static deleteBtn: ElementFinder;

  private static async deleteSelected(): Promise<void> {
    await this.deleteBtn.click();

    await retryCall(() => {
      return browser.switchTo().alert().accept();
    }, 10);
  }

  private static async openItem(column: number, value: string): Promise<void> {
    const clientSideScript = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      return row.find('a')[0];
    `;

    const itemLink = await browser.element(by.js(clientSideScript));
    const isPresent = await itemLink.isPresent();

    if (isPresent) {
      await itemLink.click();
    }
  }

  private static async selecteItem(column: number, value: string): Promise<boolean> {
    const clientSideScript = `
      const row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      return row.find('input[type="checkbox"]');
    `;
    const itemChkBox = await browser.element(by.js(clientSideScript));
    const isPresent = await itemChkBox.isPresent();
    if (isPresent) {
      await itemChkBox.click();
      return true;
    }
    return false;
  }
}
