import IHttpService = angular.IHttpService;
import IPromise = angular.IPromise;
import ICookiesService = angular.cookies.ICookiesService;

export class UserService {

  private userServiceUrl = 'rest/user';
  private authServiceUrl = 'rest/auth';

  private currentUser: UserDTO = null;


  constructor(
    private $http: IHttpService,
    private $cookies: ICookiesService) {
    
      const cookieName = 'Session-Key';
      const sessionKey = $cookies.get(cookieName);
      if(angular.isDefined(sessionKey)) {
        this.getUserFromSessionData(sessionKey);
      };
      
    }

  createUser(user: UserDTO): IPromise<UserDTO> {
    return this.$http.post(this.userServiceUrl, user)
      .then((res) => (<UserDTO>res.data));
  }

  createUserSession(user: BaseUserDTO) {
    return this.$http.post(`${this.authServiceUrl}/login`, user)
      .then((res) => {
        this.currentUser = <UserDTO>res.data 
        return (this.currentUser)
      });
  }

  getUserFromSessionData(sessionKey: string) {
    this.$http.get(`${this.userServiceUrl}`)
      .then((res) => this.currentUser = <UserDTO>res.data );
  }
  
  logout() {
    return this.$http.post(`${this.authServiceUrl}/logout`, undefined);
  }
  
  getLoggedinUser() {
    return this.currentUser;
  }
  
  isAdmin() {
    return angular.isObject(this.currentUser) && !!this.currentUser.isAdmin;
  }
}

angular.module('NavitasFitness').service('userService', UserService);

export interface BaseUserDTO {
  email: string;
  password: string;
}

export interface UserDTO extends BaseUserDTO{
  navitasId: string;
  isAdmin ?: boolean;
}
