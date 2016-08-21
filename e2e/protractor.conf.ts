


import {Config} from 'protractor';

export const config: Config = {
  framework: 'jasmine',
  baseUrl: 'http://localhost:9000',
  capabilities: {
    browserName: 'chrome'
  // }, {
  //   browserName: 'firefox'
  },
  specs: [ 'specs/*.js' ]
};