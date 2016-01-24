/// <reference path=".../../../typings/angularjs/angular"/>
/// <reference path=".../../../typings/angular-material/angular-material"/>
/// <reference path=".../../../typings/angular-ui-router/angular-ui-router"/>

import "./PageComponents/MainPage/MainPage"
import "./PageComponents/Blog/Blog"
import { RegistrationForm } from "./PageComponents/RegistrationForm/RegistrationForm"

export class AppComponent {

  constructor(
    private $mdDialog: angular.material.IDialogService,
    private $mdMedia: angular.material.IMedia) {

  }

  openRegistrationDialog(event) {
    this.$mdDialog.show({
      controller: RegistrationForm,
      // controller: DialogController,
      templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
      parent: angular.element(document.body),
      targetEvent: event,
      clickOutsideToClose: true,
      fullscreen: false
    })
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