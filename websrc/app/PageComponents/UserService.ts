import IHttpService = angular.IHttpService;
import IPromise = angular.IPromise;
import ICookiesService = angular.cookies.ICookiesService;

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUser: IUserDTO = null;

  constructor(
    private $http: IHttpService,
    private $cookies: ICookiesService) {

    const cookieName = 'Session-Key';
    const sessionKey = $cookies.get(cookieName);
    if (angular.isDefined(sessionKey)) {
      this.getUserFromSessionData();
    }
  }

  public createUser(user: IUserDTO): IPromise<IUserDTO> {
    return this.$http.post(this.userServiceUrl, user)
      .then((res) => (<IUserDTO> res.data));
  }

  public createUserSession(user: IBaseUserDTO) {
    return this.$http.post(`${this.authServiceUrl}/login`, user)
      .then((res) => {
        this.currentUser = <IUserDTO> res.data;
        return (this.currentUser);
      });
  }

  public getUserFromSessionData() {
    this.$http.get(this.userServiceUrl)
      .then((res) => {
        this.currentUser = <IUserDTO> res.data;
      });
  }

  public logout() {
    return this.$http.post(`${this.authServiceUrl}/logout`, undefined).then(() => {
      this.currentUser = null;
    });
  }

  public getLoggedinUser() {
    return this.currentUser;
  }

  public isAdmin() {
    return angular.isObject(this.currentUser) && !!this.currentUser.isAdmin;
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
