import { BlogPostsService } from './PageComponents/Blog/BlogPostsService';
import { UserService } from './PageComponents/UserService';

import { AppComponent } from './AppComponent';
import { CkEditor } from './Components/CkEditor/CkEditor';
import { BlogComponent } from './PageComponents/Blog/Blog';
import { MainPageComponent } from './PageComponents/MainPage/MainPage';
import { nfResetOnChange } from './PageComponents/RegistrationForm/nfEmailAvailable';
import { UserStatusComponent } from './PageComponents/UserStatus/UserStatus';

import IStateProvider = angular.ui.IStateProvider;
import IUrlRouterProvider = angular.ui.IUrlRouterProvider;

import ILocationProvider = angular.ILocationProvider;

import IThemingProvider = angular.material.IThemingProvider;

export const NavitasFitnessModule = angular.module('NavitasFitness', [
  'ngMaterial', 'ui.router', 'ngCookies', 'ngMessages',
  ])
  .config(($mdThemingProvider: IThemingProvider) => {
    $mdThemingProvider.theme('default')
      .primaryPalette('blue')
      .accentPalette('orange');
  })
  .config((
    $stateProvider: IStateProvider,
    $urlRouterProvider: IUrlRouterProvider,
    $locationProvider: ILocationProvider
  ) => {
    $locationProvider.html5Mode(true);
    $urlRouterProvider.otherwise("/");

    $stateProvider
      .state('MainPage', {
        template: "<main-page></main-page>",
        url: "/main-page",
      })
      .state('Blog', {
        template: "<blog></blog>",
        url: "/blog",
      })
      .state('Status', {
        template: "<user-status></user-status>",
        url: "/status",
      });
  })
  .service('userService', UserService)
  .service('blogPostsService', BlogPostsService)

  .component('ckEditor', CkEditor)
  .component('blog', BlogComponent)
  .component('mainPage', MainPageComponent)
  .component('userStatus', UserStatusComponent)
  .component('appComponent', AppComponent)

  .directive(nfResetOnChange.name, nfResetOnChange.factory);
