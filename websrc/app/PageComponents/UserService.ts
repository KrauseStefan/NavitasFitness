import { isDefined } from 'angular';
import { Observable, Subject } from 'rxjs';

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUserSubject = new Subject<IUserDTO>();

  constructor(
    private $http: ng.IHttpService,
    private $cookies: ng.cookies.ICookiesService) {

    const cookieName = 'Session-Key';
    const sessionKey = $cookies.get(cookieName);
    if (isDefined(sessionKey)) {
      this.getUserFromSessionData();
    }
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
    return this.$http.get(this.userServiceUrl)
      .then((res) => {
        const currentUser = <IUserDTO> res.data;
        this.currentUserSubject.next(currentUser);
        return currentUser;
      });
  }
}

export interface IBaseUserDTO {
  email: string;
  password: string;
}

export interface IUserDTO extends IBaseUserDTO {
  navitasId: string;
  isAdmin?: boolean;
}
