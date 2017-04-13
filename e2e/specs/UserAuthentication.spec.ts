import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { RegistrationDialogPageObject } from '../PageObjects/RegistrationDialogPageObject';
import { verifyBrowserLog } from '../utility';
import { browser, protractor } from 'protractor';

const userInfo = {
  name: 'test',
  email: 'email@domain.com',
  accessId: 'AccessId1',
  password: 'Password123',
};

const alternateUserInfo = {
  email: 'otherEmail@domain.dk',
  accessId: 'AccessId2',
};

describe('User Autentication', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] ensure test user does not exist', () => {
    browser.get('/');
    new DataStoreManipulator()
      .removeUserByEmail(userInfo.email)
      .removeUserByAccessId(userInfo.accessId)
      .destroy();
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
      accessId: userInfo.accessId,
      password: userInfo.password,
    });

    loginDialog.loginButton.click();

    expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    expect(loginDialog.errorCredentialsInvalid.isDisplayed()).toBe(true);

    loginDialog.safeClick(loginDialog.cancelButton);
  });

  describe('user creation', () => {
    let regDialog: RegistrationDialogPageObject;

    let keyUserDifferntEmail;
    let keyUserDifferntAccessId;
    let keyUser;

    beforeEach(() => {
      regDialog = NavigationPageObject.openRegistrationDialog();
    });

    afterEach(() => {
      regDialog.safeClick(regDialog.cancelButton);
    });

    it('should be able to create a user', () => {
      regDialog.fillForm({
        name: userInfo.name,
        email: alternateUserInfo.email,
        accessId: alternateUserInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      regDialog.buttonRegister.click();

      const dataStoreManipulator = new DataStoreManipulator();
      keyUserDifferntEmail = dataStoreManipulator.getUserEntityIdFromEmail(alternateUserInfo.email);
      dataStoreManipulator.destroy();

      expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should be able to override unregistred user with a differnt email but same Access ID', () => {
      regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: alternateUserInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      regDialog.buttonRegister.click();

      const dataStoreManipulator = new DataStoreManipulator();
      keyUserDifferntAccessId = dataStoreManipulator.getUserEntityIdFromEmail(userInfo.email);
      dataStoreManipulator.destroy();

      expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should be able to override unregistred user with a differnt Access ID but same email', () => {
      regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: userInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      regDialog.buttonRegister.click();

      const dataStoreManipulator = new DataStoreManipulator();
      keyUser = dataStoreManipulator.getUserEntityIdFromEmail(userInfo.email);
      dataStoreManipulator.destroy();

      expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should not be able to verify overriden user with different email', () => {
      DataStoreManipulator.sendValidationRequestFromKey(keyUserDifferntEmail).then(
        () => { fail(); },
        () => { /* success */ }
      );
    });

    it('should not be able to verify overriden user with different accessId', () => {
      DataStoreManipulator.sendValidationRequestFromKey(keyUserDifferntAccessId).then(
        () => { fail(); },
        () => { /* success */ }
      );

    });

    it('should be able to verify user email from link', () => {
      DataStoreManipulator.sendValidationRequestFromKey(keyUser);
    });

  });

  describe('registration validation', () => {
    let regDialog: RegistrationDialogPageObject;

    beforeEach(() => {
      regDialog = NavigationPageObject.openRegistrationDialog();

      regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: userInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
    });

    afterEach(() => {
      regDialog.safeClick(regDialog.cancelButton);
    });

    it('should get an error message when using an already registred email', () => {
      regDialog.buttonRegister.click();

      expect(regDialog.errorEmailConflict.isDisplayed()).toBe(true);
      verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should get an error message when using an already registred accessId', () => {
      regDialog.fillForm({ email: 'email_other@domain.com' }); // To avoid using an already registred email

      regDialog.buttonRegister.click();

      expect(regDialog.errorAccessIdConflict.isDisplayed()).toBe(true);
      verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should get an error message when using an invalid access id', () => {
      regDialog.fillForm({
        email: 'email_other@domain.com',
        accessId: 'Invalid Id',
      });

      regDialog.buttonRegister.click();

      expect(regDialog.errorAccessIdInvalid.isDisplayed()).toBe(true);
      verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should validate some user input client side', () => {
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

      regDialog.fillForm({ accessId: '' });
      expect(regDialog.buttonRegister.isEnabled()).toBe(false);

      regDialog.fillForm({ accessId: userInfo.accessId });
      expect(regDialog.buttonRegister.isEnabled()).toBe(true);
    });
  });

  it('should be able to login a user', () => {
    const loginDialog = NavigationPageObject.openLoginDialog();

    loginDialog.fillForm({
      accessId: userInfo.accessId,
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

    browser.wait(NavigationPageObject.menuLogin.isDisplayed(), 1000, 'menu did not display login button');

    expect(NavigationPageObject.menuLogin.isDisplayed()).toBe(true);
    expect(NavigationPageObject.menuRegister.isDisplayed()).toBe(true);
    expect(NavigationPageObject.menuLogout.isDisplayed()).toBe(false);

    NavigationPageObject.closeMenu();
  });

});
