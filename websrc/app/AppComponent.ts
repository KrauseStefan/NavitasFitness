import './PageComponents/MainPage/MainPage';
import './PageComponents/Blog/Blog';
import './PageComponents/UserStatus/UserStatus';
import { RegistrationForm } from './PageComponents/RegistrationForm/RegistrationForm';
import { LoginForm } from './PageComponents/LoginForm/LoginForm';
import { UserService } from './PageComponents/UserService';

import IDialogService = angular.material.IDialogService;
import IMedia = angular.material.IMedia;
import IStateProvider = angular.ui.IStateProvider;
import IUrlRouterProvider = angular.ui.IUrlRouterProvider;
import ILocationProvider = angular.ILocationProvider;

export class AppComponent {

  constructor(
    private $mdDialog: IDialogService,
    private $mdMedia: IMedia,
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
  });