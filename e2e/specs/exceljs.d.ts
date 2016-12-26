import { Writable } from "stream";

export declare namespace excel {

  type earchCb = (cell: ICell, row: number) => void;

  enum CellValueType {
    Null = 0,
    Merge = 1,
    Number = 2,
    String = 3,
    Date = 4,
    Hyperlink = 5,
    Formula = 6,
  }

  interface IAllignment {
    horizontal: 'general' | string;
    indent: number;
    shrinkToFit: boolean;
    textRotation: 0;
    vertical: 'bottom' | string;
    wrapText: boolean;
  }

  interface IFullAddress {
    address: string;
    col: 2;
    row: 1;
    sheetName: string;
  }

  interface IBorderStyle {
    style: 'thin' |
    'dotted' |
    'dashDot' |
    'hair' |
    'dashDotDot' |
    'slantDashDot' |
    'mediumDashed' |
    'mediumDashDotDot' |
    'mediumDashDot' |
    'medium' |
    'double' |
    'thick';
    up: boolean;
    down: boolean;
    color: {
      argb: string
    };
  }

  interface IBorder {
    top: IBorderStyle;
    left: IBorderStyle;
    bottom: IBorderStyle;
    right: IBorderStyle;
    diagonal: IBorderStyle;
  }

  interface IFront {
    name: string;
    family: number;
    size: number;
    underline: string;
    bold: boolean;
  }

  interface IStyle {
    alignment: IAllignment;
  }

  interface ICell {
    isMerged: boolean;
    master: IColumn;
    isHyperlink: boolean;
    hyperlink?: any;
    value: any;
    text: string;
    fullAddress: IFullAddress;
    name?: any;
    names: any[];
    style: IStyle;
    worksheet: IWorksheet;
    workbook: IWorkbook;
    dataValidation?: any;
    model: any;
    destroy: () => any;
    numFmt?: string;
    font?: IFront;
    alignment: IAllignment;
    border?: IBorder;
    fill?: any; // TODO
    address: string;
    row: number;
    col: number;
    type: CellValueType;
    effectiveType: CellValueType;
    toCsvString(): string;
    merge(ref: any): any;
    unmerge(): any;
    isMergedTo(ref: any): any;
    addName(name: any): any;
    removeName(name: any): any;
    removeAllNames(): any;
  }

  interface IEachCellOpt {
    includeEmpty: true;
  }

  interface IColumn extends String {
    header: string | string[];
    key: string;
    width: number;
    hidden: boolean;
    outlineLevel: number;
    collapsed: boolean;
    eachCell(cb: earchCb): void;
    eachCell(opt: IEachCellOpt, cb: earchCb): void;
  }

  interface IWorksheet {
    dataValidations: any;
    id: number;
    name: string;
    pageSetup: any;
    properties: any;
    views: any[];
    actualColumnCount: number;
    columnCount: number;
    actualRowCount: number;
    rowCount: number;
    columns: IColumn[];
    hasMerges: boolean;
    lasRow: any;
    model: any;
    tabColor?: any;
    workbook: IWorkbook;

    addRow(value: any): any;
    addRows(value: any): any;
    destroy(): any;
    eachRow(cb: (row: any, index: number) => void): void;
    findCell(row: number, col: number): ICell;
    findCell(ref: string): ICell;
    findRow(ref: string | number): any;
    getColumn(id: number | string): IColumn;
    getCell(row: number, col: number): ICell;
    getCell(ref: string): ICell;
    getRow(ref: string | number): any;
    getSheetValues(): any[][];
    mergeCells(...args: any[]): any;
    spliceColumns(start: any, count: any): any;
    spliceRows(start: any, count: any): any;
    unMergeCells(...args: any[]): any;
  }

  interface IWorkbook {
    new (): IWorkbook;
    dimensions: any;
    xlsx: {
      createInputStream(): Writable
    };

    getWorksheet(id: number | string): IWorksheet;

  }

}
