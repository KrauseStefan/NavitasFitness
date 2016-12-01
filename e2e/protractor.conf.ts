import { Config, ProtractorBy } from 'protractor';

const timeoutMils = 1000 * 60 * 10;

declare const angular: any;
declare const by: ProtractorBy;

export const config: Config = {
    baseUrl: 'http://localhost:8080',
    // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
    // directConnect: true,
    framework: 'jasmine',
    jasmineNodeOpts: { defaultTimeoutInterval: timeoutMils },
    multiCapabilities: [{
        browserName: 'chrome',
        maxInstances: 3,
        // shardTestFiles: true,
        // }, {
        // browserName: 'firefox',
        // marionette: true,
        // // maxInstances: 3,
        // // shardTestFiles: true,
    }],
    onPrepare: () => {
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
