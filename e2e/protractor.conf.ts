import { Config, ProtractorBy, browser } from 'protractor';

import { PluginConfig } from 'protractor/built/plugins';

const timeoutMils = 1000 * 60 * 10;

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
    disableHTMLReport: false,
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

export const config: Config = {
    baseUrl: 'http://localhost:8080',
    directConnect: true,
    // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
    // directConnect: true,
    framework: 'jasmine',
    jasmineNodeOpts: { defaultTimeoutInterval: timeoutMils },
    multiCapabilities: [{
        browserName: 'chrome',
        maxInstances: 5,
        shardTestFiles: true,
        // }, {
        // browserName: 'firefox',
        // marionette: true,
        // // maxInstances: 3,
        // shardTestFiles: true,
    }],
    onPrepare: () => {

        const disableNgAnimate = () => {
            angular.module('disableNgAnimate', []).run(['$animate', ($animate) => {
                $animate.enabled(false);
            }]);
        };
        browser.addMockModule('disableNgAnimate', disableNgAnimate);

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
    plugins: [utilsPlugin],
    // seleniumArgs: [
    // '-Dwebdriver.gecko.driver=./node_modules/protractor/node_modules/webdriver-manager/selenium/geckodriver-v0.11.1',
    // ],
    specs: ['specs/*.js'],
};
