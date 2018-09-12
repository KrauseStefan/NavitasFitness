import { module } from 'angular';

import { MainPageService } from './PageComponents/MainPage/MainPageService';
import { UserService } from './PageComponents/UserService';

import { AppComponent } from './AppComponent';
import { CkEditor } from './Components/CkEditor/CkEditor';
import { AdminPageComponent, adminRouterState } from './PageComponents/AdminPage/AdminPageComponent';
import { MainPageComponent, mainPageRouterState } from './PageComponents/MainPage/MainPage';
import { nfResetOnChange } from './PageComponents/RegistrationForm/nfResetOnChange';
import { nfShouldEqual } from './PageComponents/RegistrationForm/nfShouldEqual';
import { statusRouterState, UserStatusComponent } from './PageComponents/UserStatus/UserStatus';

import { DefaultErrorHandlingModule } from './DefaultErrorHandling';

import 'angular-animate';
import 'angular-aria';
import 'angular-cookies';
import 'angular-material';
import 'angular-messages';
import 'angular-ui-grid';
import 'angular-ui-router';

import ngMat = ng.material;

export const NavitasFitnessModule = module('NavitasFitness', [
  'ngMaterial', 'ngCookies', 'ngMessages',
  'ui.router',
  'ui.grid', 'ui.grid.selection',
  DefaultErrorHandlingModule.name,
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

  .component('ckEditor', CkEditor)
  .component('mainPage', MainPageComponent)
  .component('userStatus', UserStatusComponent)
  .component('adminPage', AdminPageComponent)
  .component('appComponent', AppComponent)

  .directive(nfResetOnChange.name, nfResetOnChange.factory)
  .directive(nfShouldEqual.name, nfShouldEqual.factory);
