import {Component, NgFor} from "angular2/angular2"
import {CkeEditorComponent} from "../ckEditor/CkEditor"


@Component({
  selector: 'main',
  templateUrl: '/main/main.html',
  directives: [CkeEditorComponent, NgFor]
})
export class Main {

  public entries: string[] = [];

  constructor() {
    for(let i = 0; i < 2; i++){
      this.entries.push(`tekst streng $i`)
    }
  }


}

