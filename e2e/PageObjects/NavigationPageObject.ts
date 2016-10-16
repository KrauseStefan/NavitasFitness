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
}
