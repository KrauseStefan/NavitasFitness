import { $, ElementFinder, browser, by, element } from 'protractor';

function byModel(model: string) {
  return element(by.model(model));
}

export class StatusPageObject {

  public static statusMsgField = byModel('$ctrl.statusMessages[$ctrl.model.statusMsgKey]');
  public static subscriptionEndField = byModel('$ctrl.model.validUntill');
  public static paypalSimBtn = $('form[action="http://localhost:8081/processPayment"] input[name="submit"]');
  public static paymentHistoryCompleatedEntry = $('tr td:nth-child(3)');

  public static getFirstRowCell(index: number): ElementFinder {
    return $(`tr td:nth-child(${index})`);
  }

  public static waitForPaypalSimBtn() {
    const btnIsDisplayed = () => StatusPageObject.paypalSimBtn.isDisplayed();
    browser.wait(btnIsDisplayed, 5 * 1000, 'Paypall button did not display in time');
  }

  public static triggerPaypalPayment() {
    browser.ignoreSynchronization = true;
    StatusPageObject.paypalSimBtn.click();
    $('a').click();

    browser.wait(browser.executeScript(() => document.readyState), 1000, 'Page did not load');

    browser.ignoreSynchronization = false;
  }

}
