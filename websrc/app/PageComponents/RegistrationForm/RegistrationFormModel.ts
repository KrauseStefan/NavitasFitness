import { UserDTO } from '../UserService';

export class RegistrationFormModel implements UserDTO {
  public name = '';
  public email = '';
  public password = '';
  public passwordRepeat = '';
  public accessId = '';
}
