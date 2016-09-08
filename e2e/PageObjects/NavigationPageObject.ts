import { by, element } from 'protractor/globals';
// import {browser, by, element, $, $$, ExpectedConditions, protractor} from 'protractor/globals';

export class NavigationPageObject {
  static allTabs = element(by.tagName('md-tabs'))

  static mainPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="MainPage"]'));
  static blogPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="Blog"]'));
  static statusPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="Status"]'));

}