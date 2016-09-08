import IScope = angular.IScope;
import IAugmentedJQuery = angular.IAugmentedJQuery;

export class CkEditorComponent {

  public content: string;
  public isEditable: boolean = false;
  public unsubscribe: Function = angular.noop;
  public editor: CKEDITOR.editor = null;

  constructor(private $scope: IScope, private $element: IAugmentedJQuery) {
    this.$scope.$watch('$ctrl.isEditable', () => {
      this.isEditable ? this.enableEditor() : this.disableEditor();
    });
  }

  public enableEditor() {
    // this.editor = CKEDITOR.replace(<any>this.getEditordiv());
    this.unsubscribe();
    this.unsubscribe = angular.noop;

    this.getEditordiv().contentEditable = 'true';

    this.editor = CKEDITOR.inline(<any> this.getEditordiv());
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

  public getEditordiv() {
    return <HTMLDivElement> this.$element[0].querySelector('.editorContent');
  }

  public updateContent(content: string) {
    this.getEditordiv().innerHTML = content;
  }

  public resetEditor() {
    this.editor.resetDirty();
  }

}

export const CkEditor = {
  bindings: {
    content: '=',
    isEditable: '=',
  },
  controller: CkEditorComponent,
  template: '<div class="editorContent"></div>',
};
