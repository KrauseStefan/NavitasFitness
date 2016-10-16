import { LoginForm } from './PageComponents/LoginForm/LoginForm';
import { RegistrationForm } from './PageComponents/RegistrationForm/RegistrationFormCtrl';
import { UserService } from './PageComponents/UserService';

import IDialogService = angular.material.IDialogService;
import IMedia = angular.material.IMedia;

class AppComponentController {

  constructor(
    private $mdDialog: IDialogService,
    private $mdMedia: IMedia,
    private userService: UserService) {
  }

  public openRegistrationDialog(event: MouseEvent) {

    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: RegistrationForm,
      fullscreen: false,
      parent: angular.element(document.body),
      targetEvent: event,
      templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
    });

  }

  public openLoginDialog(event: MouseEvent) {

    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: LoginForm,
      fullscreen: false,
      parent: angular.element(document.body),
      targetEvent: event,
      templateUrl: '/PageComponents/LoginForm/LoginForm.html',
    });

  }

  public logout() {
    this.userService.logout();
  }

  public isLoggedIn() {
    return angular.isObject(this.userService.getLoggedinUser());
  }

}

export const AppComponent = {
  controller: AppComponentController,
  templateUrl: './AppComponent.html',
};
