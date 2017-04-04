import { waitForPageToLoad } from '../utility';
import { $, browser, by, element } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

export enum TransactionTableCells {
  Amount = 1,
  PaymentDate,
  Status,
}

export interface IParsedDate { day: number; month: number; year: number; }

export function dateParts(date: string): IParsedDate {
  const seperator = '.';
  // format 'DD-MM-YYYY'
  const [day, month, year] = date
    .split(seperator)
    .map((i) => parseInt(i, 10));
  return { day, month, year };
}

function byModel(model: string) {
  return element(by.model(model));
}

function getModelValue(model: string): wdp.Promise<string> {
  return <any>byModel(model).evaluate(model);
}

export class StatusPageObject {

  public static paypalBtn = $('form[action] input[name="submit"]');

  public static termsAcceptedChkBx = $('[name="termsAccepted"]');

  public static getStatusMsgFieldValue(): wdp.Promise<string> {
    return <any>byModel('$ctrl.statusMessages[$ctrl.model.statusMsgKey]')
      .evaluate('$ctrl.model.statusMsgKey');
  }

  public static getValidUntilFieldValue(): wdp.Promise<string> {
    return getModelValue(this.subscriptionEndFieldModel);
  }

  public static getTableCellText(row: number, cell: TransactionTableCells): wdp.Promise<string> {
    if (row === 0) {
      throw "0 is an invalid index";
    }
    if (row > 0) {
      return $(`tr:nth-child(${row}) td:nth-child(${cell})`).getText();
    } else {
      return $(`tr:nth-last-child(${row * -1}) td:nth-child(${cell})`).getText();
    }
  }

  public static getFirstTransactionDate(): wdp.Promise<string> {
    return this.getTableCellText(1, TransactionTableCells.PaymentDate);
  }

  public static waitForPaypalSimBtn() {
    const btnIsDisplayed = () => StatusPageObject.paypalBtn.isDisplayed();
    browser.wait(btnIsDisplayed, 5 * 1000, 'Paypall button did not display in time');
  }

  public static triggerPaypalPayment() {
    browser.ignoreSynchronization = true;
    StatusPageObject.paypalBtn.click();

    waitForPageToLoad();
    $('a').click();

    waitForPageToLoad();
    browser.ignoreSynchronization = false;
  }

  public static getPageDates() {
    return wdp.all([
      this.getFirstTransactionDate(),
      this.getValidUntilFieldValue(),
    ]).then((values: string[]) => {
      const [firstTrxDate, validUntil] = values.map(dateParts);

      return {
        firstTrxDate,
        validUntil,
      };
    });
  }

  private static subscriptionEndFieldModel = '$ctrl.model.validUntill';
}
