import { browser } from 'protractor/globals';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
// import {ProtractorBrowser} from 'protractor';
// import {browser, by, By, element, $, $$, ExpectedConditions, protractor} from 'protractor/globals';

interface BrowserLog {
  level: {
    name: string, //SERVERE
    value: number
  },
  message: string,
  timestamp: number,
  type: string
}

describe('Navigation tests', () => {
  browser.get('/')

  afterEach(() => {
    (<any>browser).manage().logs().get('browser').then((browserLogs: BrowserLog[]) => {
      if (browserLogs.length > 0) {
        throw `Error was thrown doring test execution: [${browserLogs[0].level.name}] ${browserLogs[0].message}`
      }
    });
  });

  it('should respond to the basic "/" address', () => {
    NavigationPageObject.blogPageTab.click();
    expect(browser.getLocationAbsUrl()).toBe('/blog');

    NavigationPageObject.statusPageTab.click();
    expect(browser.getLocationAbsUrl()).toBe('/status');

    NavigationPageObject.mainPageTab.click();
    expect(browser.getLocationAbsUrl()).toBe('/main-page');
  });

  // it('should respond to the basic "/" address', () => {
  //   browser.get('/')
  // });


});