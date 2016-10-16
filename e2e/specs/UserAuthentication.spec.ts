import { verifyBrowserLog } from '../utility';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { RegistrationDialogPageObject } from '../PageObjects/RegistrationDialogPageObject';
import { LoginDialogPageObject } from '../PageObjects/LoginDialogPageObject';

const userInfo = {
  email: '20-email@domain.com',
  password: 'Password123',
  navitasId: '1234509876'
}

describe('User Autentication', () => {

  afterEach(() => verifyBrowserLog());

  it('should not be able to login before user has been created', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogin.click();
    const loginDialog = new LoginDialogPageObject();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password
    });

    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    expect(loginDialog.errorLoginSuccessful.isDisplayed()).toBe(true);

    loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to create a user', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuRegister.click();
    const regDialog = new RegistrationDialogPageObject();

    regDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
      navitasId: userInfo.navitasId
    });
    regDialog.buttonRegister.click();

    expect(regDialog.formContainer.isPresent()).toBe(false);

    regDialog.safeClick(regDialog.cancelButton);
  });

  it('should not be able to create a user that already exists', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuRegister.click();
    const regDialog = new RegistrationDialogPageObject();

    regDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password,
      passwordRepeat: userInfo.password,
      navitasId: userInfo.navitasId
    });

    regDialog.buttonRegister.click();

    expect(regDialog.formContainer.isDisplayed()).toBe(true);
    expect(regDialog.errorEmailUnavailable.isDisplayed()).toBe(true);
    verifyBrowserLog(['http://localhost:8080/rest/user 0:0 Failed to load resource: the server responded with a status of 409 (Conflict)']);

    regDialog.safeClick(regDialog.cancelButton);
  });

  it('should be able to login a user', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogin.click();
    expect(NavigationPageObject.menuLogout.isDisplayed()).toBe(false);

    const loginDialog = new LoginDialogPageObject();

    loginDialog.fillForm({
      email: userInfo.email,
      password: userInfo.password
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
  });

});