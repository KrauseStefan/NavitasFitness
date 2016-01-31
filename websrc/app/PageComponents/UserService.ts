/// <reference path=".../../../../typings/angularjs/angular.d.ts"/>

export class UserService {

  private serviceUrl = 'rest/user'

  constructor(private $http: angular.IHttpService) { }

  createUser(user: UserDTO): angular.IPromise<UserDTO> {
    return this.$http.post(this.serviceUrl, user)
      .then((res) => (<UserDTO>res.data));
  }

  createUserSession(user: BaseUserDTO) {
    return this.$http.post(`rest/auth/login`, user)
      .then((res) => (<UserDTO>res.data));
  }
  
}

angular.module('NavitasFitness').service('userService', UserService);

export interface BaseUserDTO {
  email: string;
  password: string;
}

export interface UserDTO extends BaseUserDTO{
  navitasId: string;
}
