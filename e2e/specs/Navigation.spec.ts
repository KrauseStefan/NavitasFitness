import { browser } from 'protractor';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility'

describe('Navigation tests', () => {
  browser.get('/')

  afterEach(() => verifyBrowserLog());

  it('should respond to the basic "/" address', () => {
    NavigationPageObject.blogPageTab.click();
    expect(browser.getLocationAbsUrl()).toBe('/blog');

    // NavigationPageObject.statusPageTab.click();
    // expect(browser.getLocationAbsUrl()).toBe('/status');

    NavigationPageObject.mainPageTab.click();
    expect(browser.getLocationAbsUrl()).toBe('/main-page');
  });

  // it('should respond to the basic "/" address', () => {
  //   browser.get('/')
  // });


});