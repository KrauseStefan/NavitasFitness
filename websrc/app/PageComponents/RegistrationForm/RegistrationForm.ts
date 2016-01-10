import { Component, ElementRef, AfterViewInit, AfterContentInit, OnInit } from "angular2/core"
import { Router, Location} from "angular2/router"

@Component({
  templateUrl: '/PageComponents/RegistrationForm/RegistrationForm.html',
  selector: 'registration-form'
})
export class RegistrationForm implements AfterViewInit {

  constructor(private elementRef: ElementRef, private location: Location, private router: Router) {
  }

  ngAfterViewInit() {

    //work around
    // https://github.com/PolymerElements/paper-dialog/issues/80
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

  cancel() {
    //hack until aux routes gets fixed
    const base = this.location.path().split(/[\/()]/g).filter(i => i !== '')[0]
    this.router.navigateByUrl(`/${base}`);

    this.getDialogElement().close();
  }

}