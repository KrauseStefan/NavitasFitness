import { Config } from 'protractor';

export const config: Config = {
  framework: 'jasmine',
  baseUrl: 'http://localhost:8080',
  // baseUrl: 'http://navitas-fitness-aarhus.appspot.com/',
  multiCapabilities: [{
    browserName: 'chrome'
  // }, {
  //   browserName: 'firefox'
  }],
  specs: ['specs/*.js']
};