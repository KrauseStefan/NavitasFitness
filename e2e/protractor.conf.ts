import { Config } from 'protractor';

export const config: Config = {
  baseUrl: 'http://localhost:8080',
  framework: 'jasmine',
  // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
  multiCapabilities: [{
    browserName: 'chrome',
  // }, {
  //   browserName: 'firefox'
  }],
  specs: ['specs/*.js'],
};
