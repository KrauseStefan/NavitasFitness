import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';

describe('Navigation tests', () => {

  afterEach(() => verifyBrowserLog());

  it('should respond to the basic "/" address', async () => {
    await browser.get('/');

    await NavigationPageObject.mainPageTab.click();
    const currentUrl = await browser.getCurrentUrl();
    await expect(currentUrl.endsWith('/main-page/')).toBe(true);
  });

});
