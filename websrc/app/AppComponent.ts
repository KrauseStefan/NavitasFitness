/// <reference path=".../../../typings/angularjs/angular"/>
/// <reference path=".../../../typings/angularjs/angular-cookies"/>
/// <reference path=".../../../typings/angular-material/angular-material"/>
/// <reference path=".../../../typings/angular-ui-router/angular-ui-router"/>

import "./PageComponents/MainPage/MainPage"
import "./PageComponents/Blog/Blog"
import { RegistrationForm } from "./PageComponents/RegistrationForm/RegistrationForm"
import { LoginForm } from "./PageComponents/LoginForm/LoginForm"
import { UserService } from "./PageComponents/UserService"

export class AppComponent {

  constructor(
    private $mdDialog: angular.material.IDialogService,
    private $mdMedia: angular.material.IMedia,
    private userService: UserService) {
      
  }

  openRegistrationDialog(event: MouseEvent) {

    this.$mdDialog.show({
      templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
      targetEvent: event,
      controller: RegistrationForm,
      parent: angular.element(document.body),
      clickOutsideToClose: true,
      fullscreen: false
    });

  }
  
  openLoginDialog(event: MouseEvent) {

    this.$mdDialog.show({
      templateUrl: '/PageComponents/LoginForm/LoginForm.html',
      targetEvent: event,
      controller: LoginForm,
      parent: angular.element(document.body),
      clickOutsideToClose: true,
      fullscreen: false
    });

  }
  
  logout() {
    this.userService.logout().then(() => {
      //TODO display a message
    });
  }
  
  isLoggedIn() {
    return angular.isObject(this.userService.getLoggedinUser());
  }
  
}

angular.module('NavitasFitness')
  .component('appComponent', {
    templateUrl: './AppComponent.html',
    controller: AppComponent
  })
  .config((
      $stateProvider: angular.ui.IStateProvider,
      $urlRouterProvider: angular.ui.IUrlRouterProvider
    ) => {
    $urlRouterProvider.otherwise("/")

    $stateProvider
      .state('MainPage', {
        url: "/main-page",
        template: "<main-page></main-page>",
      })
      .state('Blog', {
        url: "/blog",
        template: "<blog></blog>",
      });
  });