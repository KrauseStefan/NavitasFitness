// tslint:disable-next-line no-implicit-dependencies
import { selection } from 'ui-grid';

export class AdminUiGridConstants {

  public readonly options: Readonly<selection.IGridOptions & uiGrid.IGridOptionsOf<any>> = {
    data: [],
    enableColumnMenus: false,
    enableFiltering: true,
    enableHorizontalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    enableVerticalScrollbar: this.uiGridConstants.scrollbars.WHEN_NEEDED,
    rowHeight: 42,
    enableRowHeaderSelection: false,
  };

  constructor(private uiGridConstants: uiGrid.IUiGridConstants) {

    const headerHeight = 33;
    const filterHeight = 28;
    (<any>this.options).headerRowHeight = Math.ceil((filterHeight + headerHeight) / 2);

  }
}
