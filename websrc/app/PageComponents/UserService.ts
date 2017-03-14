import { Observable, ReplaySubject } from 'rxjs';

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUserSubject = new ReplaySubject<IUserDTO>(1);

  constructor(
    private $http: ng.IHttpService,
    private $cookies: ng.cookies.ICookiesService) {

    this.getUserFromSessionData();
  }

  public createUser(user: IUserDTO): ng.IPromise<IUserDTO> {
    return this.$http.post<IUserDTO>(this.userServiceUrl, user)
      .then((res) => res.data);
  }

  public createUserSession(user: IBaseUserDTO): ng.IPromise<IUserDTO> {
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

  public getLoggedinUser$(): Observable<IUserDTO> {
    return this.currentUserSubject.asObservable();
  }

  private getUserFromSessionData(): ng.IPromise<IUserDTO> {
    return this.$http.get<IUserSessionDto>(this.userServiceUrl)
      .then((res) => {
        const currentUser = res.data.user;
        currentUser.isAdmin = res.data.isAdmin;
        this.currentUserSubject.next(currentUser);
        return currentUser;
      }).catch(() => {
        this.currentUserSubject.next(null);
      });
  }
}

export interface IBaseUserDTO {
  email: string;
  password: string;
}

export interface IUserDTO extends IBaseUserDTO {
  name: string;
  accessId: string;
  isAdmin?: boolean;
}

export interface IUserSessionDto {
  user: IUserDTO;
  isAdmin: boolean;
}
