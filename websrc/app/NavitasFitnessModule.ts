import { module } from 'angular';

import { MainPageService } from './PageComponents/MainPage/MainPageService';
import { UserService } from './PageComponents/UserService';

import { appComponent } from './AppComponent';
import { ckEditor } from './Components/CkEditor/CkEditor';
import { adminPageModule, adminRouterState } from './PageComponents/AdminPage/AdminPageModule';
import { mainPageComponent, mainPageRouterState } from './PageComponents/MainPage/MainPage';
import { nfResetOnChange } from './PageComponents/RegistrationForm/nfResetOnChange';
import { nfShouldEqual } from './PageComponents/RegistrationForm/nfShouldEqual';
import { statusRouterState, userStatusComponent } from './PageComponents/UserStatus/UserStatus';

import { defaultErrorHandlingModule } from './DefaultErrorHandling';

import 'angular-animate';
import 'angular-aria';
import 'angular-cookies';
import 'angular-material';
import 'angular-messages';
import 'angular-ui-router';

import ngMat = ng.material;

export const navitasFitnessModule = module('NavitasFitness', [
  'ngMaterial', 'ngCookies', 'ngMessages',
  'ui.router',
  defaultErrorHandlingModule.name,
  adminPageModule.name,
])
  .config(($mdThemingProvider: ngMat.IThemingProvider) => {
    $mdThemingProvider.theme('default')
      .primaryPalette('blue')
      .accentPalette('orange');
  })

  .config((
    $stateProvider: ng.ui.IStateProvider,
    $urlRouterProvider: ng.ui.IUrlRouterProvider,
    $locationProvider: ng.ILocationProvider,
    $sceDelegateProvider: ng.ISCEDelegateProvider,
  ) => {
    $locationProvider.html5Mode(true);
    $urlRouterProvider.otherwise('/main-page/');

    $stateProvider
      .state('MainPage', mainPageRouterState)
      .state('Status', statusRouterState)
      .state('Admin', adminRouterState);

    $sceDelegateProvider.resourceUrlWhitelist([
      // Allow same origin resource loads.
      'self',
      'http://localhost:8081/processPayment',
      'https://www.sandbox.paypal.com/cgi-bin/webscr',
      'https://www.paypal.com/cgi-bin/webscr',
    ]);

  })

  .run(($window: ng.IWindowService, $q: ng.IQService) => {
    $window['Promise'] = $q;
  })

  .service('userService', UserService)
  .service('mainPageService', MainPageService)

  .component('ckEditor', ckEditor)
  .component('mainPage', mainPageComponent)
  .component('userStatus', userStatusComponent)
  .component('appComponent', appComponent)

  .directive(nfResetOnChange.name, nfResetOnChange.factory)
  .directive(nfShouldEqual.name, nfShouldEqual.factory);
