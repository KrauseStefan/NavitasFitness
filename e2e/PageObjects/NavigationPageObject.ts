import { LoginDialogPageObject } from './LoginDialogPageObject';
import { RegistrationDialogPageObject } from './RegistrationDialogPageObject';
import { $ } from 'protractor';

export class NavigationPageObject {
  public static allTabs = $('md-tabs');
  public static menuContent = $('md-menu-content');

  public static menuButton = $('.md-button[ng-click="$mdOpenMenu()"]');

  public static mainPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="MainPage"]');
  public static blogPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="Blog"]');
  public static statusPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="Status"]');

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
