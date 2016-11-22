import { $ } from 'protractor';

export class FrontPageObject {
  public static adminEditBtn = $('button[ng-click="$ctrl.entry.enabled = !$ctrl.entry.enabled"]');

  public static adminSaveBtn = $('button[ng-click="$ctrl.saveEntry($ctrl.entry)"]');

  public static editableArea = $('ck-editor .editorContent');
};
