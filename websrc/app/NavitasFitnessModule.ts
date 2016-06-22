import { UserService } from './PageComponents/UserService';
import { BlogPostsService } from './PageComponents/Blog/BlogPostsService';

import { AppComponent } from './AppComponent';
import { CkEditor } from './Components/CkEditor/CkEditor';
import { BlogComponent } from './PageComponents/Blog/Blog';
import { MainPageComponent } from './PageComponents/MainPage/MainPage';
import { UserStatusComponent } from './PageComponents/UserStatus/UserStatus';

import IStateProvider = angular.ui.IStateProvider;
import IUrlRouterProvider = angular.ui.IUrlRouterProvider;
import ILocationProvider = angular.ILocationProvider;
import IThemingProvider = angular.material.IThemingProvider;


export const NavitasFitnessModule = angular.module('NavitasFitness', ['ngMaterial', 'ui.router', 'ngCookies'])
  .config(function ($mdThemingProvider: IThemingProvider) {
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
        url: "/main-page",
        template: "<main-page></main-page>",
      })
      .state('Blog', {
        url: "/blog",
        template: "<blog></blog>",
      })
      .state('Status', {
        url: "/status",
        template: "<user-status></user-status>",
      });
  })
  .service('userService', UserService)
  .service('blogPostsService', BlogPostsService)

  .component('ckEditor', CkEditor)
  .component('blog', BlogComponent)
  .component('mainPage', MainPageComponent)
  .component('userStatus', UserStatusComponent)
  .component('appComponent', AppComponent);
