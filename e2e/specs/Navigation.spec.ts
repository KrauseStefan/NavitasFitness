import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { browser } from 'protractor';

describe('Navigation tests', () => {

  afterEach(() => verifyBrowserLog());

  it('should respond to the basic "/" address', () => {
    browser.get('/');
    // NavigationPageObject.statusPageTab.click();
    // expect(browser.getLocationAbsUrl()).toBe('/status');

    NavigationPageObject.mainPageTab.click();
    expect(
      browser.getCurrentUrl().then((absUrl) => absUrl.endsWith('/main-page/'))
    ).toBe(true);
  });

  // it('should respond to the basic "/" address', () => {
  //   browser.get('/')
  // });
});
