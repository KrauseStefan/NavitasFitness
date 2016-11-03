import { Config } from 'protractor';

const timeoutMils = 1000 * 60 * 10;

export const config: Config = {
  baseUrl: 'http://localhost:8080',
  // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
  framework: 'jasmine',
  jasmineNodeOpts: {defaultTimeoutInterval: timeoutMils},
  multiCapabilities: [{
    browserName: 'chrome',
  // }, {
  //   browserName: 'firefox'
  }],
  specs: ['specs/*.js'],
};
