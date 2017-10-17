import { LoginForm } from './PageComponents/LoginForm/LoginForm';
import { KeyFieldName, PasswordChangeFormCtrl } from './PageComponents/PasswordChangeForm/PasswordChangeFormCtrl';
import { RegistrationForm } from './PageComponents/RegistrationForm/RegistrationFormCtrl';
import { IUserDTO, UserService } from './PageComponents/UserService';
import { element, isDefined, isObject } from 'angular';

class AppComponentController {

  public tabs = [{
    id: 'MainPage',
    icon: 'home',
    text: 'Home',
    disabled: () => false,
  }, {
    id: 'Status',
    icon: 'perm_identity',
    text: 'Payment Status',
    disabled: () => !this.isLoggedIn(),
  }];

  public selectedTabIndex: number;
  private loggedInUser: IUserDTO;

  constructor(
    private $mdDialog: ng.material.IDialogService,
    private $mdMedia: ng.material.IMedia,
    private $location: ng.ILocationService,
    private $scope: ng.IScope,
    private $state: ng.ui.IStateService,
    private userService: UserService) {

    this.updateSelectedTab();

    const searchParams = this.$location.search();
    if (isDefined(searchParams[KeyFieldName])) {
      this.openPasswordChangeDialog();
    }

    userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      this.loggedInUser = user;
      this.updateSelectedTab();
    });
  }

  public updateSelectedTab() {
    let index = this.tabs.findIndex(i => i.id === this.$state.current.name);
    if (!this.tabs[index].disabled()) {
      this.$scope.$applyAsync(() => {
        // Selection cannot be changed while the tab is disabled, it takes a digest cycle before it is unlocked
        this.selectedTabIndex = index;
      });
    }
  }

  public openPasswordChangeDialog() {

    this.$mdDialog.show({
      clickOutsideToClose: true,
      controller: PasswordChangeFormCtrl,
      fullscreen: false,
      parent: element(document.body),
      templateUrl: '/PageComponents/PasswordChangeForm/PasswordChangeForm.html',
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
    this.userService.logout().then(() => {
      this.loggedInUser = null;
    });
  }

  public isLoggedIn() {
    return isObject(this.loggedInUser);
  }

}

export const AppComponent = {
  controller: AppComponentController,
  templateUrl: './AppComponent.html',
};
