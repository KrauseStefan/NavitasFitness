/// <reference path=".../../../../../typings/angularjs/angular.d.ts"/>

export class UserService {

  private serviceUrl = 'rest/user'

  constructor(private $http: angular.IHttpService) { }

  createUser(user: UserDTO): angular.IPromise<UserDTO> {
    console.log(user);

    return this.$http.post(this.serviceUrl, user)
      .then((res) => (<UserDTO>res.data));
  }

  createUserSession(user: UserDTO) {
  }
}

angular.module('NavitasFitness').service('userService', UserService);

export class UserDTO {
  email: string
  password: string
  navitasId: string
}