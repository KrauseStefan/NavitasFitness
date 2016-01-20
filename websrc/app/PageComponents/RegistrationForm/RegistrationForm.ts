import { Component, ElementRef, AfterViewInit, AfterContentInit, AfterViewChecked, OnInit } from "angular2/core"
import { Router, Location} from "angular2/router"
// import "rxjs/observable/"

import { UserService, UserDTO } from './UserService'

export class RegistrationFormModel {
  email: string = ""
  emailRepeat: string = ""
  password: string = ""
  passwordRepeat: string = ""
  navitasId: string = ""

  toUserDTO(): UserDTO {
    return {
      email: this.email,
      password: this.password,
      navitasId: this.navitasId
    }
  }
}

@Component({
  templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
  selector: 'registration-form'
})
export class RegistrationForm implements AfterViewInit {

  submitted = false;
  model = new RegistrationFormModel();

  constructor(
    private elementRef: ElementRef,
    private location: Location,
    private router: Router,
    private userService: UserService) {
  }

  ngAfterViewInit() {
    // work around
    // https://github.com/PolymerElements/paper-dialog-scrollable/issues/13

    window.setTimeout(() => {
      if(!this.getDialogElement().opened){
        this.getDialogElement().open();
      }
      window.setTimeout(() => {
        this.getDialogElement().fit();
      });
    });
  }

  getDialogElement() {
    return this.elementRef.nativeElement.getElementsByTagName('paper-dialog')[0];
  }

  submit() {
    this.userService.createUser(this.model.toUserDTO()).subscribe(()=>{}, ()=>{}, ()=>{
      this.model = new RegistrationFormModel();
      this.close();
    });
  }

  close() {
    //hack until aux routes gets fixed
    const base = this.location.path().split(/[\/()]/g).filter(i => i !== '')[0]
    this.router.navigateByUrl(`/${base}`);

    this.getDialogElement().close();
  }

  getDebugModel() {
    JSON.stringify(this.model, null, 2)
  }
}