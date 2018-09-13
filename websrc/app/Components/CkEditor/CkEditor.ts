import { IComponentOptions, noop } from 'angular';

import 'ckeditor';

export class CkEditorComponent {

  public content = '';
  public isEditable: boolean = false;
  public unsubscribe = noop;
  public editor: CKEDITOR.editor | null = null;

  constructor(private $scope: ng.IScope, private $element: ng.IAugmentedJQuery) {
    this.$scope.$watch('$ctrl.isEditable', () => {
      this.isEditable ? this.enableEditor() : this.disableEditor();
    });
  }

  public enableEditor() {
    this.unsubscribe();
    this.unsubscribe = noop;

    this.getEditordiv().contentEditable = 'true';

    this.editor = CKEDITOR.inline(<any>this.getEditordiv(), { customConfig: '' });
    this.editor.on('change', (event) => {
      this.$scope.$apply(() => this.content = event.editor.getData());
    });
  }

  public disableEditor() {
    if (this.editor !== null) {
      this.editor.destroy();
      this.editor = null;
      this.getEditordiv().contentEditable = 'false';
    }

    this.unsubscribe = this.$scope.$watch('$ctrl.content', ((content: string) => {
      this.updateContent(content);
    }));
  }

  public getEditordiv(): HTMLElement {
    const div = (<HTMLElement>(<any>this.$element)[0]).querySelector('.editorContent');

    if (div && div instanceof HTMLElement) {
      return div;
    }

    throw new Error('Unable to find Ckeditor div');
  }

  public updateContent(content: string) {
    this.getEditordiv().innerHTML = content;
  }

  public resetEditor() {
    if (this.editor) {
      this.editor.resetDirty();
    }
  }

}

export const ckEditor: IComponentOptions = {
  bindings: {
    content: '=',
    isEditable: '=',
  },
  controller: CkEditorComponent,
  template: '<div class="editorContent"></div>',
};
