import { AlerDialogPageObject } from '../PageObjects/AlertDialogPageObject';
import { DataStoreManipulator } from '../PageObjects/DataStoreManipulator';
import { NavigationPageObject } from '../PageObjects/NavigationPageObject';
import { RegistrationDialogPageObject } from '../PageObjects/RegistrationDialogPageObject';

import { browser, protractor } from 'protractor';
import { verifyBrowserLog } from '../utility';

const userInfo = {
  name: 'test',
  email: 'auth-test-1@domain.com',
  accessId: 'auth-test-1',
  password: 'Password123',
};

const alternateUserInfo = {
  email: 'auth-test-2@domain.dk',
  accessId: 'auth-test-2',
};

describe('User Autentication', () => {

  afterEach(() => verifyBrowserLog());

  it('[META] ensure test user does not exist', async () => {
    await browser.get('/');
    await DataStoreManipulator.loadUserKinds();
    await DataStoreManipulator.removeUserByEmail(userInfo.email);
    await DataStoreManipulator.removeUserByAccessId(userInfo.accessId);
    await DataStoreManipulator.removeUserByEmail(alternateUserInfo.email);
    await DataStoreManipulator.removeUserByAccessId(alternateUserInfo.accessId);
  });

  it('[META] ensure user is not logged in', async () => {
    await NavigationPageObject.menuButton.click();
    const isDisplayed = await NavigationPageObject.menuLogout.isDisplayed();
    if (isDisplayed) {
      await NavigationPageObject.menuLogout.click();
    } else {
      await NavigationPageObject.menuLogin.sendKeys(protractor.Key.ESCAPE);
    }
  });

  it('should not be able to login before user has been created', async () => {
    const loginDialog = await NavigationPageObject.openLoginDialog();

    await loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });

    await loginDialog.loginButton.click();

    await expect(loginDialog.formContainer.isDisplayed()).toBe(true);
    await expect(loginDialog.errorCredentialsInvalid.isDisplayed()).toBe(true);

    await loginDialog.safeClick(loginDialog.cancelButton);
  });

  describe('user creation', () => {
    let regDialog: RegistrationDialogPageObject;

    let keyUserDifferntEmail;
    let keyUserDifferntAccessId;
    let keyUser;

    beforeEach(async () => {
      regDialog = await NavigationPageObject.openRegistrationDialog();
    });

    afterEach(async () => {
      await regDialog.safeClick(regDialog.cancelButton);
    });

    it('should be able to create a user', async () => {
      await regDialog.fillForm({
        name: userInfo.name,
        email: alternateUserInfo.email,
        accessId: alternateUserInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      await regDialog.buttonRegister.click();
      await AlerDialogPageObject.mainButton.click();

      await DataStoreManipulator.loadUserKinds();
      keyUserDifferntEmail = await DataStoreManipulator.getUserEntityIdFromEmail(alternateUserInfo.email);

      await expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should be able to override unregistred user with a differnt email but same Access ID', async () => {
      await regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: alternateUserInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      await regDialog.buttonRegister.click();
      await AlerDialogPageObject.mainButton.click();

      await DataStoreManipulator.loadUserKinds();
      keyUserDifferntAccessId = await DataStoreManipulator.getUserEntityIdFromEmail(userInfo.email);

      await expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should be able to override unregistred user with a differnt Access ID but same email', async () => {
      await regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: userInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
      await regDialog.buttonRegister.click();
      await AlerDialogPageObject.mainButton.click();

      await DataStoreManipulator.loadUserKinds();
      keyUser = await DataStoreManipulator.getUserEntityIdFromEmail(userInfo.email);

      await expect(regDialog.formContainer.isPresent()).toBe(false);
    });

    it('should not be able to verify overriden user with different email', async () => {
      await DataStoreManipulator.sendValidationRequestFromKey(keyUserDifferntEmail).then(
        () => { fail(); },
        () => { /* success */ },
      );
    });

    it('should not be able to verify overriden user with different accessId', async () => {
      await DataStoreManipulator.sendValidationRequestFromKey(keyUserDifferntAccessId).then(
        () => { fail(); },
        () => { /* success */ },
      );
    });

    it('should be able to verify user email from link', async () => {
      await DataStoreManipulator.sendValidationRequestFromKey(keyUser);
    });

  });

  describe('registration validation', async () => {
    let regDialog: RegistrationDialogPageObject;

    beforeEach(async () => {
      regDialog = await NavigationPageObject.openRegistrationDialog();

      await regDialog.fillForm({
        name: userInfo.name,
        email: userInfo.email,
        accessId: userInfo.accessId,
        password: userInfo.password,
        passwordRepeat: userInfo.password,
      });
    });

    afterEach(async () => {
      await regDialog.safeClick(regDialog.cancelButton);
    });

    it('should get an error message when using an already registred email', async () => {
      await regDialog.buttonRegister.click();

      await expect(regDialog.errorEmailConflict.isDisplayed()).toBe(true);
      await verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should get an error message when using an already registred accessId', async () => {
      await regDialog.fillForm({ email: 'email_other@domain.com' }); // To avoid using an already registred email

      await regDialog.buttonRegister.click();

      await expect(regDialog.errorAccessIdConflict.isDisplayed()).toBe(true);
      await verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should get an error message when using an invalid access id', async () => {
      await regDialog.fillForm({
        email: 'email_other@domain.com',
        accessId: 'Invalid Id',
      });

      await regDialog.buttonRegister.click();

      await expect(regDialog.errorAccessIdInvalid.isDisplayed()).toBe(true);
      await verifyBrowserLog([[
        'http://localhost:8080/rest/user',
        '0:0',
        'Failed to load resource: the server responded with a status of 409 (Conflict)',
      ].join(' ')]);
    });

    it('should validate some user input client side', async () => {
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);
      await expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

      await regDialog.fillForm({ name: '' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);

      await regDialog.fillForm({ name: userInfo.name });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);

      await regDialog.fillForm({ password: 'bad' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);
      await expect(regDialog.errorPasswordDifferent.isPresent()).toBe(true);

      await regDialog.fillForm({ passwordRepeat: 'bad' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);
      await expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

      await regDialog.fillForm({ passwordRepeat: userInfo.password });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);
      await expect(regDialog.errorPasswordDifferent.isPresent()).toBe(true);

      await regDialog.fillForm({ password: userInfo.password });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);

      await expect(regDialog.errorPasswordDifferent.isPresent()).toBe(false);

      await regDialog.fillForm({ email: 'ab' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);

      await regDialog.fillForm({ email: '' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);

      await regDialog.fillForm({ email: userInfo.email });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);

      await regDialog.fillForm({ accessId: '' });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(false);

      await regDialog.fillForm({ accessId: userInfo.accessId });
      await expect(regDialog.buttonRegister.isEnabled()).toBe(true);
    });
  });

  it('should be able to login a user', async () => {
    const loginDialog = await NavigationPageObject.openLoginDialog();

    await loginDialog.fillForm({
      accessId: userInfo.accessId,
      password: userInfo.password,
    });

    await loginDialog.loginButton.click();

    await expect(loginDialog.formContainer.isPresent()).toBe(false);

    await loginDialog.safeClick(loginDialog.cancelButton);
  });

  it('should be able to logout current user', async () => {
    await NavigationPageObject.menuButton.click();

    await expect(NavigationPageObject.menuLogin.isDisplayed()).toBe(false);
    await expect(NavigationPageObject.menuRegister.isDisplayed()).toBe(false);

    await NavigationPageObject.menuLogout.click();
    await NavigationPageObject.menuButton.click();

    await browser.wait(NavigationPageObject.menuLogin.isDisplayed(), 1000, 'menu did not display login button');

    await expect(NavigationPageObject.menuRegister.isDisplayed()).toBe(true);
    await expect(NavigationPageObject.menuLogout.isDisplayed()).toBe(false);

    await NavigationPageObject.closeMenu();
  });

});
