import { Component, ElementRef, AfterViewInit, AfterContentInit, AfterViewChecked, OnInit } from "angular2/core"
import {NgForm}    from 'angular2/common';
import { Router, Location} from "angular2/router"

export class RegistrationFormModel {
  email: string = ""
  emailRepeat: string = ""
  password: string = ""
  passwordRepeat: string = ""
  navitasId: string = ""
}

@Component({
  templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
  selector: 'registration-form'
})
export class RegistrationForm implements AfterViewInit {

  submitted = false;
  model = new RegistrationFormModel(); //= {
  //   email: 'test@mail.com',
  //   emailRepeat: 'test@mail.com',
  //   password: '1234567',
  //   passwordRepeat: '1234567',
  //   navitasId: '1234567'
  // }

  constructor(private elementRef: ElementRef, private location: Location, private router: Router) {
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

  onSubmit() {
    this.submitted = true;
  }

  cancel() {
    //hack until aux routes gets fixed
    const base = this.location.path().split(/[\/()]/g).filter(i => i !== '')[0]
    this.router.navigateByUrl(`/${base}`);

    this.getDialogElement().close();
  }

}