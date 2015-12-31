import { provide, Injectable } from "angular2/core"
import { Http, Response }  from "angular2/http"



@Injectable()
export class UserServices {

  constructor(private http: Http) { }

  createUser(user: UserDTO) {

  }

  createUserSession(user: UserDTO) {

  }

}

export class UserDTO {
  email: string
  password: string
  navitasId: string
}