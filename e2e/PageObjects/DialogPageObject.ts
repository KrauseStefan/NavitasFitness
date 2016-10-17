import { $, ElementFinder } from 'protractor';
import { promise as wdpromise } from 'selenium-webdriver';

export type formNameValuesMap = { [name: string]: string }

export class DialogPageObject {
  public formContainer = $('md-dialog');

  public fillForm(formValues: formNameValuesMap) {
    Object.keys(formValues).forEach(name => {
      const field = this.formContainer.$(`input[name="${name}"]`);
      field.clear();
      field.sendKeys(formValues[name]);
    });
  }

  public safeClick(element: ElementFinder): wdpromise.Promise<any> {
    const resolved = wdpromise.fullyResolved(null);
    return element.isDisplayed().then<any>((isDisplayed) => {
      return !isDisplayed ? resolved : element.click();
    }, () => {/* */});
  }
}
