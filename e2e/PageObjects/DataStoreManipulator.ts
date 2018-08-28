import { retryCall, waitForPageToLoad } from '../utility';
import { DataStoreClientScripts } from './DataStoreClientScripts';

import * as http from 'http';
import { ElementFinder, ProtractorBrowser, browser as mainBrowser } from 'protractor';

let browser: ProtractorBrowser;
let clientScriptsProxy: DataStoreClientScripts;

export class DataStoreManipulator {

  public static async loadUserKinds(): Promise<void> {
    await browser.waitForAngularEnabled(false);
    const kind = 'User';
    await browser.get(`http://localhost:8000/datastore?kind=${kind}`);
    await browser.ready;

    await browser.waitForAngularEnabled(false);
    await browser.executeScript(`window.clientScripts = new ${DataStoreClientScripts.toString()}`);

    this.deleteBtn = await browser.$('#delete_button');
  }

  public static async init(): Promise<void> {
    browser = await mainBrowser.forkNewDriverInstance(false, false, false);
    clientScriptsProxy = DataStoreClientScripts.getProxy(browser);
  }

  public static async destroy() {
    await waitForPageToLoad();
    await browser.close();
    await browser.quit();
  }

  public static async performEmailVerification(email: string): Promise<void> {
    const key = await DataStoreManipulator.getUserEntityIdFromEmail(email);
    await browser.sleep(500); // below call might fail silently without a delay here
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

        reject('User email could not be verified using URL: ' + url);
      });
    });
  }

  public static async getUserEntityIdFromEmail(email: string): Promise<string> {
    const key = await clientScriptsProxy.getValue('email', email, 'key');
    if (key) {
      return key;
    }

    throw `Unable to lookup user DB key, email used: ${email}`;
  }

  public static async getUserEntityResetSecretFromEmail(email: string): Promise<string> {
    const secret = await clientScriptsProxy.getValue('email', email, 'PasswordResetSecret');
    if (secret) {
      return secret;
    }

    throw `Unable to lookup reset secret, email used: ${email}`;
  }

  public static async removeUserByAccessId(accessId: string): Promise<void> {
    try {
      const checkbox = await clientScriptsProxy.getRowCheckbox('AccessId', accessId);
      await checkbox.click();

      await DataStoreManipulator.deleteSelected();
    } catch (e) {
      console.log(`Failed to remove user: ${accessId}, user not found! Ignoring error`);
    }
  }

  public static async removeUserByEmail(email: string): Promise<void> {
    try {
      const checkbox = await clientScriptsProxy.getRowCheckbox('email', email);
      await checkbox.click();

      await DataStoreManipulator.deleteSelected();
    } catch (e) {
      console.log(`Failed to remove user: ${email}, user not found! Ignoring error`);
    }
  }

  public static async makeUserAdmin(email: string): Promise<void> {
    const link = await clientScriptsProxy.getRowIdLink('email', email);
    await link.click();

    const selectAdmin = `document.querySelector('select[name="bool|IsAdmin"]').value = 1;`;
    await browser.driver.executeScript(selectAdmin);
    await browser.$('input[value="Save Changes"]').click();
  }

  private static deleteBtn: ElementFinder;

  private static async deleteSelected(): Promise<void> {
    await this.deleteBtn.click();

    await retryCall(() => {
      return Promise.resolve(browser.switchTo().alert().accept());
    }, 10);
  }

}
