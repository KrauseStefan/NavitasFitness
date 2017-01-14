import { retryCall } from '../utility';
import { ElementFinder, ProtractorBrowser, browser as mainBrowser, by } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

let browser: ProtractorBrowser;

export class DataStoreManipulator {

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

  public removeUser(email: string) {
    this.selecteItem(8, email);

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

    return this;
  }

  public makeUserAdmin(email) {
    this.openItem(8, email);

    const selectAdmin = `document.querySelector('select[name="bool|IsAdmin"]').value = 1;`;
    browser.driver.executeScript(selectAdmin);
    browser.$('input[value="Save Changes"]').click();
    return this;
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
