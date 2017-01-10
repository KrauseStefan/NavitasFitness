import { retryCall } from '../utility';
import { $, ElementFinder } from 'protractor';
import { promise as wdp } from 'selenium-webdriver';

export type formNameValuesMap = { [name: string]: string }

export class DialogPageObject {

  private static fillField(field: ElementFinder, value): wdp.Promise<string> {
    field.clear();
    field.sendKeys(value);

    return field.getAttribute('value').then(text => {
      if (value === text) {
        return wdp.fullyResolved(value);
      }

      return wdp.rejected(`expected value: "${value}" did not equal field value: "${text}"`);
    });
  }

  public formContainer = $('md-dialog');

  public fillForm(formValues: formNameValuesMap) {
    Object.keys(formValues).forEach(name => {
      const field = this.formContainer.$(`input[name="${name}"]`);
      // Sometimes form fields fail to sendKeys (one might be missing)
      // https://github.com/angular/protractor/issues/698
      retryCall(() => DialogPageObject.fillField(field, formValues[name]), 3);
    });
  }

  public safeClick(element: ElementFinder): wdp.Promise<any> {
    const resolved = wdp.fullyResolved(null);
    return element.isDisplayed().then<any>((isDisplayed) => {
      return !isDisplayed ? resolved : element.click();
    }, () => {/* */ });
  }
}
