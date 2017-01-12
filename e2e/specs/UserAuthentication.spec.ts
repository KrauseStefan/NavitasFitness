import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { verifyBrowserLog } from '../utility';
import { browser, protractor } from 'protractor';

const userInfo = {
  name: 'test',
  email: 'email@domain.com',
  navitasId: '1234509876',
  password: 'Password123',
};

describe('User Autentication', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] ensure test user does not exist', () => {
    browser.get('/');
    new DataStoreManipulator().removeUser(userInfo.email).destroy();
  });

  it('[META] ensure user is not logged in', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogout.isDisplayed().then((isDisplayed) => {
      if (isDisplayed) {
        NavigationPageObject.menuLogout.click();
      } else {
        browser.actions().sendKeys(protractor.Key.ESCAPE).perform();
      }
    });

  });

  it('should not be able to login before user has been created', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    expect(loginDialog.errorLoginSuccessful.isDisplayed()).toBe(true);

    loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to create a user', () => {
    const regDialog = NavigationPageObject.openRegistrationDialog();

    regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      navitasId: userInfo.navitasId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.termsAcceptedChkBx.click();
    regDialog.buttonRegister.click();

    expect(regDialog.formContainer.isPresent()).toBe(false);

    regDialog.safeClick(regDialog.cancelButton);
  });

  it('should not be able to create a user that already exists', () => {
    const regDialog = NavigationPageObject.openRegistrationDialog();

    regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      navitasId: userInfo.navitasId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.termsAcceptedChkBx.click();

    regDialog.buttonRegister.click();

    expect(regDialog.formContainer.isDisplayed()).toBe(true);
    expect(regDialog.errorEmailUnavailable.isDisplayed()).toBe(true);
    verifyBrowserLog([[
      'http://localhost:8080/rest/user',
      '0:0',
      'Failed to load resource: the server responded with a status of 409 (Conflict)',
    ].join(' ')]);

    regDialog.safeClick(regDialog.cancelButton);
  });

  it('should validate user input', () => {
    const regDialog = NavigationPageObject.openRegistrationDialog();
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);
    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

    regDialog.fillForm({
      name: userInfo.name,
      email: userInfo.email,
      navitasId: userInfo.navitasId,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
    });
    regDialog.termsAcceptedChkBx.click();

    expect(regDialog.buttonRegister.isEnabled()).toBe(true);
    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

    regDialog.fillForm({ name: '' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);

    regDialog.fillForm({ name: userInfo.name });
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);

    regDialog.fillForm({ password: 'bad' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);
    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(true);

    regDialog.fillForm({ passwordRepeat: 'bad' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);
    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

    regDialog.fillForm({ passwordRepeat: userInfo.password });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);
    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(true);

    regDialog.fillForm({ password: userInfo.password });
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);

    expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

    regDialog.fillForm({ email: 'ab' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);

    regDialog.fillForm({ email: '' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);

    regDialog.fillForm({ email: userInfo.email });
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);

    regDialog.fillForm({ navitasId: '' });
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);

    regDialog.fillForm({ navitasId: userInfo.navitasId });
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);

    regDialog.termsAcceptedChkBx.click();
    expect(regDialog.buttonRegister.isEnabled()).toBe(false);

    regDialog.termsAcceptedChkBx.click();
    expect(regDialog.buttonRegister.isEnabled()).toBe(true);

    regDialog.safeClick(regDialog.cancelButton);
  });

  it('should be able to login a user', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isPresent()).toBe(false);

    loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to logout current user', () => {
    NavigationPageObject.menuButton.click();

    expect(NavigationPageObject.menuLogin.isDisplayed()).toBe(false);
    expect(NavigationPageObject.menuRegister.isDisplayed()).toBe(false);

    NavigationPageObject.menuLogout.click();

    NavigationPageObject.menuButton.click();

    expect(NavigationPageObject.menuLogin.isDisplayed()).toBe(true);
    expect(NavigationPageObject.menuRegister.isDisplayed()).toBe(true);
    expect(NavigationPageObject.menuLogout.isDisplayed()).toBe(false);

    NavigationPageObject.closeMenu();
  });

});
