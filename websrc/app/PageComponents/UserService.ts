import { Observable, ReplaySubject } from 'rxjs';

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUserSubject = new ReplaySubject<UserDTO | null>(1);

  constructor(
    private $http: ng.IHttpService,
    private $q: ng.IQService,
    private $log: ng.ILogService) {

    this.getUserFromSessionData();
  }

  public createUser(user: UserDTO): ng.IPromise<UserDTO> {
    return this.$http.post<UserDTO>(this.userServiceUrl, user)
      .then((res) => res.data);
  }

  public doUserLogin(user: BaseUserDTO): ng.IPromise<UserDTO> {
    return this.$http.post<UserDTO>(`${this.authServiceUrl}/login`, user)
      .then((res) => {
        const currentUser = res.data;
        this.currentUserSubject.next(currentUser);
        return (currentUser);
      });
  }

  public logout(): ng.IPromise<void> {
    return this.$http.post(`${this.authServiceUrl}/logout`, undefined).then(() => {
      this.currentUserSubject.next(null);
    });
  }

  public getLoggedinUser$(): Observable<UserDTO | null> {
    return this.currentUserSubject.asObservable();
  }

  public sendResetPasswordRequest(email: string) {
    return this.$http.post<string>(`${this.userServiceUrl}/resetPassword/${email}`, undefined);
  }

  public sendPasswordChangeRequest(dto: ChangePasswordDTO) {
    return this.$http.post(`${this.userServiceUrl}/changePassword`, dto);
  }

  private getUserFromSessionData(): void {
    this.$http.get<UserSessionDto>(this.userServiceUrl).then((res) => {
      return this.$q<UserDTO | null>((resolve) => {
        if (!res.data || !res.data.user) {
          this.$log.debug('No User Session');
          return resolve(null);
        }

        const currentUser = res.data.user;
        currentUser.isAdmin = res.data.isAdmin;
        resolve(currentUser);
      });
    }).then((userDto) => {
      this.currentUserSubject.next(userDto);
    }, () => {
      this.currentUserSubject.next(null);
    });
  }
}

export interface ChangePasswordDTO {
  password: string;
  key: string;
  secret: string;
}

export interface BaseUserDTO {
  accessId: string;
  password: string;
}

export interface UserDTO extends BaseUserDTO {
  email: string;
  name: string;
  isAdmin?: boolean;
}

export interface UserSessionDto {
  user: UserDTO;
  isAdmin: boolean;
}
