import { IUserDTO } from '../UserService';

export class RegistrationFormModel implements IUserDTO {
  public name = '';
  public email = '';
  public password = '';
  public passwordRepeat = '';
  public accessId = '';
}
