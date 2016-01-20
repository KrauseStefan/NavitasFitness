import { provide, Injectable } from "angular2/core"
import { Http, Response }  from "angular2/http"
import {Observable} from "rxjs/Observable"
import "rxjs/add/operator/map"

@Injectable()
export class UserService {

  private serviceUrl = 'rest/user'

  constructor(private http: Http) { }

  createUser(user: UserDTO): Observable<UserDTO> {
    console.log(user);

    const data = JSON.stringify(user);
    return this.http.post(this.serviceUrl, data)
      .map((res: Response) => (<UserDTO>res.json()))

  }

  createUserSession(user: UserDTO) {
  }
}

export class UserDTO {
  email: string
  password: string
  navitasId: string
}