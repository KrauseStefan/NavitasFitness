import { ProtractorBrowser, browser as mainBrowser, by } from 'protractor';
import { promise as wdpromise } from 'selenium-webdriver';

let browser: ProtractorBrowser;

export class DataStoreManipulator {

  constructor() {
    browser = mainBrowser.forkNewDriverInstance(false, false);
    browser.ignoreSynchronization = true;
    browser.driver.get('http://localhost:8000/datastore?kind=User');
  }

  public destroy() {
    browser.sleep(1000);
    browser.quit();
  }

  public removeUser(email: string) {
    this.selecteItem(7, email);

    browser.$('#delete_button').isEnabled().then(displayed => {
      if (displayed) {
        browser.$('#delete_button').click();
        return browser.switchTo().alert().accept();
      } else {
        return wdpromise.fullyResolved<void>({});
      }
    });

    return this;
  }

  public makeUserAdmin(email) {
    this.openItem(7, email);

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
    return itemLink.click().then(undefined, () => {
      // tslint:disable-next-line
      console.log('openItem, could not find: ', value);
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
    return itemChkBox.click().then(undefined, () => {
      // tslint:disable-next-line
      console.log('selecteItem, could not find: ', value);
    });
  }
}
