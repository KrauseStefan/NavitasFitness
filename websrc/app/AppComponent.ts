import { LoginForm } from './PageComponents/LoginForm/LoginForm';
import { RegistrationForm } from './PageComponents/RegistrationForm/RegistrationFormCtrl';
import { IUserDTO, UserService } from './PageComponents/UserService';
import { element, isObject } from 'angular';

class AppComponentController {

  private loggedInUser: IUserDTO;

  constructor(
    private $mdDialog: ng.material.IDialogService,
    private $mdMedia: ng.material.IMedia,
    private userService: UserService) {

    userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      this.loggedInUser = user;
    });

  }

  public openRegistrationDialog(event: MouseEvent) {

    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: RegistrationForm,
      fullscreen: false,
      parent: element(document.body),
      targetEvent: event,
      templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
    });

  }

  public openLoginDialog(event: MouseEvent) {

    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: LoginForm,
      fullscreen: false,
      parent: element(document.body),
      targetEvent: event,
      templateUrl: '/PageComponents/LoginForm/LoginForm.html',
    });

  }

  public logout() {
    this.userService.logout();
    this.loggedInUser = null;
  }

  public isLoggedIn() {
    return isObject(this.loggedInUser);
  }

}

export const AppComponent = {
  controller: AppComponentController,
  templateUrl: './AppComponent.html',
};
