import { module } from 'angular';

import { BlogPostsService } from './PageComponents/Blog/BlogPostsService';
import { MainPageService } from './PageComponents/MainPage/MainPageService';
import { UserService } from './PageComponents/UserService';

import { AppComponent } from './AppComponent';
import { CkEditor } from './Components/CkEditor/CkEditor';
import { BlogComponent } from './PageComponents/Blog/Blog';
import { MainPageComponent } from './PageComponents/MainPage/MainPage';
import { nfResetOnChange } from './PageComponents/RegistrationForm/nfResetOnChange';
import { nfShouldEqual } from './PageComponents/RegistrationForm/nfShouldEqual';
import { UserStatusComponent, statusRouterState } from './PageComponents/UserStatus/UserStatus';

export const NavitasFitnessModule = module('NavitasFitness', [
  'ngMaterial', 'ui.router', 'ngCookies', 'ngMessages',
  ])
  .config(($mdThemingProvider: ng.material.IThemingProvider) => {
    $mdThemingProvider.theme('default')
      .primaryPalette('blue')
      .accentPalette('orange');
  })
  .config((
    $stateProvider: ng.ui.IStateProvider,
    $urlRouterProvider: ng.ui.IUrlRouterProvider,
    $locationProvider: ng.ILocationProvider
  ) => {
    $locationProvider.html5Mode(true);
    $urlRouterProvider.otherwise('/main-page');

    $stateProvider
      .state('MainPage', {
        template: '<main-page></main-page>',
        url: '/main-page',
      })
      .state('Blog', {
        template: '<blog></blog>',
        url: '/blog',
      })
      .state('Status', statusRouterState);
  })
  .service('userService', UserService)
  .service('blogPostsService', BlogPostsService)
  .service('mainPageService', MainPageService)

  .component('ckEditor', CkEditor)
  .component('blog', BlogComponent)
  .component('mainPage', MainPageComponent)
  .component('userStatus', UserStatusComponent)
  .component('appComponent', AppComponent)

  .directive(nfResetOnChange.name, nfResetOnChange.factory)
  .directive(nfShouldEqual.name, nfShouldEqual.factory);
