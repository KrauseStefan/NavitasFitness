import IScope = angular.IScope;
import IAugmentedJQuery = angular.IAugmentedJQuery;

export class CkEditorComponent {

  content: string;
  isEditable: boolean = false;
  unsubscribe: Function = angular.noop;
  editor: CKEDITOR.editor = null;

  constructor(private $scope: IScope, private $element: IAugmentedJQuery) {
    this.$scope.$watch('$ctrl.isEditable', () => {
      this.isEditable ? this.enableEditor() : this.disableEditor();
    });
  }

  enableEditor() {
    // this.editor = CKEDITOR.replace(<any>this.getEditordiv());
    this.unsubscribe();
    this.unsubscribe = angular.noop;

    this.getEditordiv().contentEditable = 'true';

    this.editor = CKEDITOR.inline(<any>this.getEditordiv());
    this.editor.on('change', (event) => {
      this.$scope.$apply(() => this.content = event.editor.getData());
    });
  }

  disableEditor() {
    if (this.editor !== null) {
      this.editor.destroy();
      this.editor = null;
      this.getEditordiv().contentEditable = 'false';
    }

    this.unsubscribe = this.$scope.$watch('$ctrl.content', ((content: string) => {
      this.updateContent(content);
    }));
  }

  getEditordiv() {
    return <HTMLDivElement> this.$element[0].querySelector('.editorContent');
  }

  updateContent(content: string) {
    this.getEditordiv().innerHTML = content;
  }

  resetEditor() {
    this.editor.resetDirty();
  }

}
angular.module('NavitasFitness')
  .component('ckEditor', {
    template: '<div class="editorContent"></div>',
    controller: CkEditorComponent,
    bindings: {
      content: '=',
      isEditable: '='
    }
  });
