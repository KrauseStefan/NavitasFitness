import { $, browser, by, element } from 'protractor';
import { waitForPageToLoad } from '../utility';

export enum TransactionTableCells {
  Amount = 1,
  PaymentDate,
  Status,
}

export interface ParsedDate { day: number; month: number; year: number; }

export function dateParts(date: string): ParsedDate {
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

function getModelValue(model: string): Promise<string> {
  return <any>byModel(model).evaluate(model);
}

export class StatusPageObject {

  public static paypalBtn = $('form[name="PaymentSandBoxForm"] input[name="submit"]');

  public static termsAcceptedChkBx = $('[name="termsAccepted"]');

  public static getStatusMsgFieldValue(): Promise<string> {
    return <any>byModel('$ctrl.statusMessages[$ctrl.model.statusMsgKey]')
      .evaluate('$ctrl.model.statusMsgKey');
  }

  public static getValidUntilFieldValue(): Promise<string> {
    return getModelValue(this.subscriptionEndFieldModel);
  }

  public static getTableCellText(row: number, cell: TransactionTableCells): Promise<string> {
    if (row === 0) {
      throw new Error('0 is an invalid index');
    }
    if (row > 0) {
      return Promise.resolve($(`tr:nth-child(${row}) td:nth-child(${cell})`).getText());
    } else {
      return Promise.resolve($(`tr:nth-last-child(${row * -1}) td:nth-child(${cell})`).getText());
    }
  }

  public static getFirstTransactionDate(): Promise<string> {
    return this.getTableCellText(1, TransactionTableCells.PaymentDate);
  }

  public static async waitForPaypalSimBtn(): Promise<boolean> {
    const btnIsDisplayed = () => StatusPageObject.paypalBtn.isDisplayed();
    return browser.wait(btnIsDisplayed, 5 * 1000, 'Paypall button did not display in time');
  }

  public static async triggerPaypalPayment(): Promise<void> {
    await browser.waitForAngularEnabled(false);
    await StatusPageObject.paypalBtn.click();

    await waitForPageToLoad();

    await $('a').click();

    await waitForPageToLoad();

    await browser.waitForAngularEnabled(true);
  }

  public static getPageDates(): Promise<{firstTrxDate: ParsedDate, validUntil: ParsedDate }> {
    return Promise.all([
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
