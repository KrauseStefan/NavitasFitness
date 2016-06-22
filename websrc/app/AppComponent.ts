import { RegistrationForm } from './PageComponents/RegistrationForm/RegistrationForm';
import { LoginForm } from './PageComponents/LoginForm/LoginForm';
import { UserService } from './PageComponents/UserService';

import IDialogService = angular.material.IDialogService;
import IMedia = angular.material.IMedia;

class AppComponentController {

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

export const AppComponent = {
  templateUrl: './AppComponent.html',
  controller: AppComponentController
};
