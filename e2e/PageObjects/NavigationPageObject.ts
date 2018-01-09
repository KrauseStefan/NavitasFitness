import { LoginDialogPageObject } from './LoginDialogPageObject';
import { RegistrationDialogPageObject } from './RegistrationDialogPageObject';
import { $, by, element } from 'protractor';

export class NavigationPageObject {

  public static allTabs = $('md-tabs');
  public static menuContent = $('md-menu-content');

  public static menuButton = $('.md-button[ng-click="$mdMenu.open()"]');

  public static mainPageTab = element((<any>by).linkUiSref('MainPage'));
  public static statusPageTab = element((<any>by).linkUiSref('Status'));

  public static menuRegister = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openRegistrationDialog($event)"]');
  public static menuLogin = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openLoginDialog($event)"]');
  public static menuLogout = NavigationPageObject.menuContent.$('[ng-click="$ctrl.logout($event)"]');

  public static async openLoginDialog(): Promise<LoginDialogPageObject> {
    await NavigationPageObject.menuButton.click();
    await NavigationPageObject.menuLogin.click();
    return new LoginDialogPageObject();
  }

  public static async openRegistrationDialog(): Promise<RegistrationDialogPageObject> {
    await NavigationPageObject.menuButton.click();
    await NavigationPageObject.menuRegister.click();
    return new RegistrationDialogPageObject();
  }

  public static async closeMenu() {
    await $('.md-menu-backdrop').click();
  }
}
