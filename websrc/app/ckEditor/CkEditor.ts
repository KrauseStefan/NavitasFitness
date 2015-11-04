/// <reference path="../../typings/ckeditor/ckeditor"/>

import { Component, ElementRef } from "angular2/angular2"

@Component({
  selector: 'ck-editor',
  templateUrl: '/ckEditor/CkEditor.html'
})

export class CkeEditorComponent {

  constructor( elementRef: ElementRef ) {
    CKEDITOR.replace( elementRef.nativeElement );
  }

}
