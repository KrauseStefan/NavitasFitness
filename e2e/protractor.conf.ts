import * as commandLineArgs from 'command-line-args';
import { Config, ProtractorBy, browser } from 'protractor';
import { PluginConfig } from 'protractor/built/plugins';

import { DataStoreManipulator } from './PageObjects/DataStoreManipulator';

const optionDefinitions = [
  { name: 'parallel', type: Boolean, defaultOption: false },
];

const cmdOpts = commandLineArgs(optionDefinitions);

const timeoutMils = 1000 * 60;

declare const angular: any;
declare const by: ProtractorBy;

interface IJasmine2ProtractorUtilsConfig extends PluginConfig {
  clearFoldersBeforeTest?: boolean;
  disableHTMLReport?: boolean;
  disableScreenshot?: boolean;
  failTestOnErrorLog?: {
    excludeKeywords: string[], // {A JSON Array}
    failTestOnErrorLogLevel: number,
  };
  htmlReportDir?: string;
  screenshotOnExpectFailure?: boolean;
  screenshotOnSpecFailure?: boolean;
  screenshotPath?: string;
}

const utilsPlugin: IJasmine2ProtractorUtilsConfig = {
  clearFoldersBeforeTest: true,
  disableHTMLReport: true,
  disableScreenshot: false,
  failTestOnErrorLog: {
    excludeKeywords: [], // {A JSON Array}
    failTestOnErrorLogLevel: 900,
  },
  htmlReportDir: './reports/htmlReports',
  package: 'jasmine2-protractor-utils',
  screenshotOnExpectFailure: true,
  screenshotOnSpecFailure: true,
  screenshotPath: './screenshots',
};

const webdriverFolder = 'node_modules/protractor/node_modules/webdriver-manager/selenium/';

export const config: Config = {
  SELENIUM_PROMISE_MANAGER: false,
  jvmArgs: [
    // '-Dwebdriver.ie.driver=${webdriverFolder}IEDriverServer3.4.0.exe',
    `-Dwebdriver.gecko.driver=${webdriverFolder}geckodriver-v0.16.1`,
  ],
  baseUrl: 'http://localhost:8080',
  // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
  // directConnect: true,
  framework: 'jasmine2',
  jasmineNodeOpts: {
    defaultTimeoutInterval: timeoutMils,
    realtimeFailure: true,
  },
  disableChecks: true,
  allScriptsTimeout: 60000,
  multiCapabilities: [{
    browserName: 'chrome',
    maxInstances: 5,
    shardTestFiles: cmdOpts.parallel,
    chromeOptions: { args: ["--headless", "--disable-gpu", "--window-size=1920,1080"] },
    // }, {
    // browserName: 'edge',
    // maxInstances: 1,
    // shardTestFiles: cmdOpts.parallel,
    // }, {
    // browserName: 'internet explorer',
    // maxInstances: 1,
    // shardTestFiles: cmdOpts.parallel,
    // }, {
    //     browserName: 'firefox',
    //     maxInstances: 3,
    //     marionette: true,
    //     shardTestFiles: cmdOpts.parallel,
  }],
  onPrepare: async () => {
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

    await DataStoreManipulator.init();
  },
  onComplete: async () => await DataStoreManipulator.destroy(),
  plugins: [utilsPlugin],
  // seleniumArgs: [
  // '-Dwebdriver.gecko.driver=./node_modules/protractor/node_modules/webdriver-manager/selenium/geckodriver-v0.11.1',
  // ],
  specs: ['specs/*.js'],
};
