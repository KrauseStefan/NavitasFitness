/// <reference path="../../typings/ckeditor/ckeditor"/>

import { Component, ElementRef, Input, OnChanges, SimpleChange } from 'angular2/angular2'
import { BlogEntry } from '../main/BlogPostsService'

@Component({
  selector: 'ck-editor',
  templateUrl: '/ckEditor/CkEditor.html'
})
export class CkeEditorComponent implements OnChanges {

  @Input() content: string = '';
  @Input() isEditable: boolean = false;


  constructor(private elementRef: ElementRef) {
  }

  enableEditor() {
    CKEDITOR.replace(this.elementRef.nativeElement);
  }

  public onChanges(changes: {[key: string]: SimpleChange}) {
    for(const key in changes) {
      console.log(`onChanges - ${key} =`, changes[key].currentValue);
    }

    if(changes['isEditable'] && changes['isEditable'].currentValue){
      this.enableEditor();
    }
  }

}
