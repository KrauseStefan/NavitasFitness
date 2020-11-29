import * as commandLineArgs from 'command-line-args';
import { browser, Config, ProtractorBy } from 'protractor';
import * as HtmlReporter from 'protractor-beautiful-reporter';

const optionDefinitions = [
  { name: 'parallel', type: Boolean, defaultOption: false },
];

const cmdOpts = commandLineArgs(optionDefinitions);

const second = 1000;
const testTimeout = 10 * second;
const scriptTimeout = 10 * second;

declare const angular: any;
declare const by: ProtractorBy;

const webdriverFolder = 'node_modules/protractor/node_modules/webdriver-manager/selenium/';

export const config: Config = {
  SELENIUM_PROMISE_MANAGER: false,
  jvmArgs: [
    `-Dwebdriver.gecko.driver=${webdriverFolder}geckodriver-v0.16.1`,
  ],
  baseUrl: 'http://localhost:8080',
  // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
  // directConnect: true,
  framework: 'jasmine2',
  jasmineNodeOpts: {
    defaultTimeoutInterval: testTimeout,
    failFast: true,
    realtimeFailure: true,
  },
  disableChecks: true,
  allScriptsTimeout: scriptTimeout,
  multiCapabilities: [{
    browserName: 'chrome',
    maxInstances: 5,
    shardTestFiles: cmdOpts.parallel,
    chromeOptions: { args: ['--headless', '--window-size=1920,1080'] },
    // }, {
    //     browserName: 'firefox',
    //     maxInstances: 3,
    //     marionette: true,
    //     shardTestFiles: cmdOpts.parallel,
  }],
  onPrepare: async () => {
    jasmine.getEnv().addReporter(new HtmlReporter({
      baseDirectory: 'reports',
      clientDefaults: {
        columnSettings: {
          displayTime: true,
          displayBrowser: false,
          displaySessionId: false,
          displayOS: false,
          inlineScreenshots: true,
        },
        preserveDirectory: false,
        searchSettings: {
          allselected: false,
          passed: false,
          failed: true,
          pending: true,
          withLog: true,
        },
      },
    }).getJasmine2Reporter());

    function disableNgAnimate() {
      angular.module('disableNgAnimate', []).run(['$animate', ($animate) => $animate.enabled(false)]);
    }
    browser.addMockModule('disableNgAnimate', disableNgAnimate);

    function disableCssAnimate() {
      angular.module('disableCssAnimate', [])
        .run(() => {
          const style = document.createElement('style');
          style.type = 'text/css';
          style.innerHTML = '* {' +
            '-webkit-transition: none !important;' +
            '-moz-transition: none !important' +
            '-o-transition: none !important' +
            '-ms-transition: none !important' +
            'transition: none !important' +
            '}';
          document.getElementsByTagName('head')[0].appendChild(style);
        });
    }
    browser.addMockModule('disableCssAnimate', disableCssAnimate);

    by.addLocator('linkUiSref', (toState: string, optParentElement: HTMLElement) => {
      const using = optParentElement || document;
      const tabs = using.querySelectorAll('md-tab-item');
      for (let i = 0; tabs.length > i; i++) {
        const uiRef = angular.element(tabs[i]).scope().tab.element.attr('ui-sref');
        if (uiRef === toState) {
          return tabs[i];
        }
      }
      return null;
    });
  },
  // seleniumArgs: [
  // '-Dwebdriver.gecko.driver=./node_modules/protractor/node_modules/webdriver-manager/selenium/geckodriver-v0.11.1',
  // ],
  specs: ['specs/*.js'],
};
