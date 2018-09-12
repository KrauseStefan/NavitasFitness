import { $, ElementFinder } from 'protractor';
import { retryCall } from '../utility';

export interface IFormNameValuesMap {
  [name: string]: string;
}

export class DialogPageObject {

  private static async fillField(field: ElementFinder, value: string): Promise<string> {
    await field.clear();
    await field.sendKeys(value);

    return new Promise<string>((resolve, reject) => {
      field.getAttribute('value').then((text) => {
        if (value === text) {
          return resolve(value);
        }

        return reject(`expected value: "${value}" did not equal field value: "${text}"`);
      }, reject);
    });
  }

  public formContainer = $('md-dialog');

  public fillForm(formValues: IFormNameValuesMap): Promise<void> {
    const promises = Object.keys(formValues).map((name) => {
      const field = this.formContainer.$(`input[name="${name}"]`);
      // Sometimes form fields fail to sendKeys (one might be missing)
      // https://github.com/angular/protractor/issues/698
      return retryCall(() => DialogPageObject.fillField(field, formValues[name]), 3);
    });
    return (<any>Promise.all(promises));
  }

  public safeClick(element: ElementFinder): Promise<void> {
    const resolved = Promise.resolve(<void>undefined);
    return Promise.resolve(element.isDisplayed().then((isDisplayed) => {
      return !isDisplayed ? resolved : element.click();
    }, () => Promise.resolve<void>(<void>undefined)));
  }
}
