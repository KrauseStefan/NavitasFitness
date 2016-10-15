import { verifyBrowserLog } from '../utility';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { RegistrationDialogPageObject } from '../PageObjects/RegistrationDialogPageObject';
import { LoginDialogPageObject } from '../PageObjects/LoginDialogPageObject';

const userInfo = {
  email: '13-email@domain.com',
  password: 'Password123',
  navitasId: '1234509876'
}

describe('User Autentication', () => {

  afterEach(() => verifyBrowserLog());

  it('should not be able to login before user has been created', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogin.click();
    const loginDialog = new LoginDialogPageObject();

    loginDialog.fieldEmail.sendKeys(userInfo.email);
    loginDialog.fieldPassword.sendKeys(userInfo.password);
    loginDialog.buttonLogin.click();

    expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    expect(loginDialog.errorLoginSuccessful.isDisplayed()).toBe(true);

    loginDialog.buttonCancel.isDisplayed().then((isDisplayed) => {
      if (isDisplayed) {
        loginDialog.buttonCancel.click();
      }
    }, () => { });
  });

  it('should be able to create a user', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuRegister.click();
    const registrationDialog = new RegistrationDialogPageObject();

    registrationDialog.fieldEmail.sendKeys(userInfo.email);
    registrationDialog.fieldPassword.sendKeys(userInfo.password);
    registrationDialog.fieldPasswordRepeat.sendKeys(userInfo.password);
    registrationDialog.fieldNavitasId.sendKeys(userInfo.navitasId);
    registrationDialog.buttonRegister.click();

    expect(registrationDialog.formContainer.isPresent()).toBe(false);

    registrationDialog.buttonCancel.isPresent().then((isDisplayed) => {
      if (isDisplayed) {
        registrationDialog.buttonCancel.click();
      }
    }, () => { });
  });

  it('should not be able to create a user that already exists', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuRegister.click();
    const registrationDialog = new RegistrationDialogPageObject();

    registrationDialog.fieldEmail.sendKeys(userInfo.email);
    registrationDialog.fieldPassword.sendKeys(userInfo.password);
    registrationDialog.fieldPasswordRepeat.sendKeys(userInfo.password);
    registrationDialog.fieldNavitasId.sendKeys(userInfo.navitasId);
    registrationDialog.buttonRegister.click();

    expect(registrationDialog.formContainer.isDisplayed()).toBe(true);
    expect(registrationDialog.errorEmailUnavailable.isDisplayed()).toBe(true);
    verifyBrowserLog(['http://localhost:8080/rest/user 0:0 Failed to load resource: the server responded with a status of 409 (Conflict)']);

    registrationDialog.buttonCancel.isDisplayed().then((isDisplayed) => {
      if (isDisplayed) {
        registrationDialog.buttonCancel.click();
      }
    }, () => { });
  });

  it('should be able to login a user', () => {
    NavigationPageObject.menuButton.click();
    NavigationPageObject.menuLogin.click();
    const loginDialog = new LoginDialogPageObject();

    loginDialog.fieldEmail.sendKeys(userInfo.email);
    loginDialog.fieldPassword.sendKeys(userInfo.password);
    loginDialog.buttonLogin.click();

    expect(loginDialog.formContainer.isPresent()).toBe(false);

    loginDialog.buttonCancel.isDisplayed().then((isDisplayed) => {
      if (isDisplayed) {
        loginDialog.buttonCancel.click();
      }
    }, () => { });
  });

});