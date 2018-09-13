import { Observable, ReplaySubject } from 'rxjs';

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUserSubject = new ReplaySubject<IUserDTO | null>(1);

  constructor(
    private $http: ng.IHttpService,
    private $q: ng.IQService,
    private $log: ng.ILogService) {

    this.getUserFromSessionData();
  }

  public createUser(user: IUserDTO): ng.IPromise<IUserDTO> {
    return this.$http.post<IUserDTO>(this.userServiceUrl, user)
      .then((res) => res.data);
  }

  public doUserLogin(user: IBaseUserDTO): ng.IPromise<IUserDTO> {
    return this.$http.post<IUserDTO>(`${this.authServiceUrl}/login`, user)
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

  public getLoggedinUser$(): Observable<IUserDTO | null> {
    return this.currentUserSubject.asObservable();
  }

  public sendResetPasswordRequest(email: string) {
    return this.$http.post<string>(`${this.userServiceUrl}/resetPassword/${email}`, undefined);
  }

  public sendPasswordChangeRequest(dto: IChangePasswordDTO) {
    return this.$http.post(`${this.userServiceUrl}/changePassword`, dto);
  }

  private getUserFromSessionData(): void {
    this.$http.get<IUserSessionDto>(this.userServiceUrl).then((res) => {
      return this.$q<IUserDTO | null>((resolve) => {
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

export interface IChangePasswordDTO {
  password: string;
  key: string;
  secret: string;
}

export interface IBaseUserDTO {
  accessId: string;
  password: string;
}

export interface IUserDTO extends IBaseUserDTO {
  email: string;
  name: string;
  isAdmin?: boolean;
}

export interface IUserSessionDto {
  user: IUserDTO;
  isAdmin: boolean;
}
