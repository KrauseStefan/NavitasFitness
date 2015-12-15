/// <reference path="../../typings/ckeditor/ckeditor"/>

import { Component, ElementRef, Input, Output, OnChanges, SimpleChange, EventEmitter } from 'angular2/core'
import { BlogEntryDTO } from '../main/BlogPostsService'

@Component({
  selector: 'ck-editor',
  templateUrl: '/ckEditor/CkEditor.html'
})
export class CkEditorComponent implements OnChanges {

  @Input() content: string = '';
  @Input() isEditable: boolean = false;

  @Output() contentChange = new EventEmitter<string>();

  editor: CKEDITOR.editor = null;

  constructor(private elementRef: ElementRef) { }

  enableEditor() {
    // this.editor = CKEDITOR.replace(<any>this.getEditordiv());
    this.getEditordiv().contentEditable = 'true';
    this.editor = CKEDITOR.inline(<any>this.getEditordiv());
    this.editor.on('change', (event) => {
      this.content = event.editor.getData();
      this.contentChange.next(this.content)
    });
  }

  disableEditor() {
    if (this.editor !== null) {
      this.editor.destroy();
      this.editor = null;
      this.getEditordiv().contentEditable = 'false';
    }
  }

  getEditordiv(): HTMLDivElement {
    return this.elementRef.nativeElement.querySelector('.editorContent');
  }

  updateContent(content) {
    this.getEditordiv().innerHTML = content;
  }

  public ngOnChanges(changes: { [key: string]: SimpleChange }) {
    // for (const key in changes) {
    //   console.log(`onChanges - ${key} =`, changes[key].currentValue);
    // }

    if (changes['content'] && !this.isEditable) {
      this.updateContent(changes['content'].currentValue);
    }

    if (changes['isEditable']) {
      changes['isEditable'].currentValue ? this.enableEditor() : this.disableEditor();
    }
  }

  resetEditor() {
    this.editor.resetDirty();
  }

}
