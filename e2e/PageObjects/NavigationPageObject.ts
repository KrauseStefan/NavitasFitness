import { by, element } from 'protractor/globals';
// import {browser, by, element, $, $$, ExpectedConditions, protractor} from 'protractor/globals';

if(element.all === undefined) {
  throw "element.all was undefined that should never be possible"
}
const allElements = element.all;

export class NavigationPageObject {
  static allTabs = allElements(by.tagName('md-tab-item'))

}