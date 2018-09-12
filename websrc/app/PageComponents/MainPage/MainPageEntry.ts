import { MainPageDTO } from './MainPageDTO';

export class MainPageEntry {
  public mainPage: MainPageDTO;
  public enabled = false;

  constructor(mainPage: MainPageDTO = new MainPageDTO()) {
    this.mainPage = mainPage;
  }
}
