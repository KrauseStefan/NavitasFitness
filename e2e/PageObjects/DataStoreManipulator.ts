import { ProtractorBrowser, browser, by } from 'protractor';

export class DataStoreManipulator {

    private browserHandle: ProtractorBrowser;

    constructor() {
        this.browserHandle = browser.forkNewDriverInstance(false, false);
        this.browserHandle.ignoreSynchronization = true;
        this.browserHandle.driver.get('http://localhost:8000/datastore');
    }

    public removeUser(email: string) {
        this.selecteItem(7, email)
        this.browserHandle.$('#delete_button').click();

        this.browserHandle.switchTo().alert().accept();
    }

    public destroy() {
        this.browserHandle.quit();
    }

    private selecteItem(column: number, value: string) {
        const getLink = `
        var row = $('.ae-table.ae-settings-block tr')
            .slice(1)
            .filter((_, elm) => $(elm).find('td:nth(${column})').text() === '${value}');

        row.find('input[type="checkbox"]').click();
        `;

        this.browserHandle.driver.executeScript(getLink);
    }
}
