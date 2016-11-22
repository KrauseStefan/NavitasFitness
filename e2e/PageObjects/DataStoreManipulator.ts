import { ProtractorBrowser, browser } from 'protractor';

export class DataStoreManipulator {

  private browserHandle: ProtractorBrowser;

  constructor() {
    this.browserHandle = browser.forkNewDriverInstance(false, false);
    this.browserHandle.ignoreSynchronization = true;
    this.browserHandle.driver.get('http://localhost:8000/datastore');
  }

  public destroy() {
    this.browserHandle.quit();
  }

  public removeUser(email: string) {
    this.selecteItem(7, email);

    this.browserHandle.$('#delete_button').isDisplayed().then(displayed => {
      if (displayed) {
        this.browserHandle.$('#delete_button').click();
        this.browserHandle.switchTo().alert().accept();
      }
    }, () => {
      // do nothing
    });
    return this;
  }

  public makeUserAdmin(email) {
    this.openItem(7, email);

    const selectAdmin = `document.querySelector('select[name="bool|IsAdmin"]').value = 1;`;
    this.browserHandle.driver.executeScript(selectAdmin);
    this.browserHandle.$('input[value="Save Changes"]').click();
    return this;
  }

  private openItem(column: number, value: string) {
    const getLink = `
      var row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      row.find('a')[0].click();
   `;

    return this.browserHandle.driver.executeScript(getLink);
  }

  private selecteItem(column: number, value: string) {
    const getLink = `
      var row = $('.ae-table.ae-settings-block tr')
        .slice(1)
        .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

      row.find('input[type="checkbox"]').click();
   `;

    return this.browserHandle.driver.executeScript(getLink);
  }
}
