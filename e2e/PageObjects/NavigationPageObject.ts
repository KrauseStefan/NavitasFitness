import { by, element } from 'protractor';

export class NavigationPageObject {
  static allTabs = element(by.tagName('md-tabs'));
  static menuContent = element(by.tagName('md-menu-content'));

  static menuButton = element(by.css('.md-button[ng-click="$mdOpenMenu()"]'));

  static mainPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="MainPage"]'));
  static blogPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="Blog"]'));
  static statusPageTab = NavigationPageObject.allTabs.element(by.css('md-tab-item [ui-sref="Status"]'));

  static menuRegister = NavigationPageObject.menuContent.element(by.css('[ng-click="$ctrl.openRegistrationDialog($event)"]'));
  static menuLogin = NavigationPageObject.menuContent.element(by.css('[ng-click="$ctrl.openLoginDialog($event)"]'));
}

export class RegistrationDialogPageObject {
  public formContainer = element(by.tagName('md-dialog'));

  public fieldEmail = this.formContainer.element(by.css('input[ng-model="model.email"]'));
  public fieldPassword = this.formContainer.element(by.css('input[ng-model="model.password"]'));
  public fieldPasswordRepeat = this.formContainer.element(by.css('input[ng-model="model.passwordRepeat"]'));
  public fieldNavitasId = this.formContainer.element(by.css('input[ng-model="model.navitasId"]'));

  public buttonRegister = this.formContainer.element(by.css('button[ng-click="submit()"]'));
  public buttonCancel = this.formContainer.element(by.css('button[ng-click="cancel()"]'));

  public errorEmailUnavailable = this.formContainer.element(by.css('[ng-message="emailAvailable"]'));
}

export class LoginDialogPageObject {
  public formContainer = element(by.tagName('md-dialog'));

  public fieldEmail = this.formContainer.element(by.css('input[ng-model="model.email"]'));
  public fieldPassword = this.formContainer.element(by.css('input[ng-model="model.password"]'));

  public buttonLogin = this.formContainer.element(by.css('button[ng-click="submit()"]'));
  public buttonCancel = this.formContainer.element(by.css('button[ng-click="cancel()"]'));

  public errorLoginSuccessful = this.formContainer.element(by.css('[ng-message="loginSuccessful"]'));
}