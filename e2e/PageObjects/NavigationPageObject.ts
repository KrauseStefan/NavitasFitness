import { $ } from 'protractor';

export class NavigationPageObject {
  static allTabs = $('md-tabs');
  static menuContent = $('md-menu-content');

  static menuButton = $('.md-button[ng-click="$mdOpenMenu()"]');

  static mainPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="MainPage"]');
  static blogPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="Blog"]');
  static statusPageTab = NavigationPageObject.allTabs.$('md-tab-item [ui-sref="Status"]');

  static menuRegister = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openRegistrationDialog($event)"]');
  static menuLogin = NavigationPageObject.menuContent.$('[ng-click="$ctrl.openLoginDialog($event)"]');
  static menuLogout = NavigationPageObject.menuContent.$('[ng-click="$ctrl.logout($event)"]');
}
