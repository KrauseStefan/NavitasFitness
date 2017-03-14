import { LoginDialogPageObject } from './LoginDialogPageObject';
import { RegistrationDialogPageObject } from './RegistrationDialogPageObject';
import { $, by, element } from 'protractor';

export class NavigationPageObject {

  public static allTabs = $('md-tabs');
  public static menuContent = $('md-menu-content');

  public static menuButton = $('.md-button[ng-click="$mdMenu.open()"]');

  public static mainPageTab = element((<any> by).linkUiSref('MainPage'));
  public static statusPageTab = element((<any> by).linkUiSref('Status'));

  public static menuRegister = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openRegistrationDialog($event)"]');
  public static menuLogin = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openLoginDialog($event)"]');
  public static menuLogout = NavigationPageObject.menuContent.$('[ng-click="$ctrl.logout($event)"]');

  public static openLoginDialog(): LoginDialogPageObject {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogin.click();
    return new LoginDialogPageObject();
  }

  public static openRegistrationDialog(): RegistrationDialogPageObject {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuRegister.click();
    return new RegistrationDialogPageObject();
  }

  public static closeMenu() {
    $('.md-menu-backdrop').click();
  }
}
