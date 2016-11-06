import {UserService} from '../UserService';
import {MainPageDTO, MainPageService} from './MainPageService';

export class MainPage {

  public entry: MainPageEntry = new MainPageEntry({
    content: '',
    date: '',
    key: null,
    lastEditedBy: '',
  });

  constructor(public mainPageService: MainPageService, private userService: UserService) {
    mainPageService.getMainPage()
      .then(mainPage => this.entry = new MainPageEntry(mainPage));
  }

  public saveEntry(entry: MainPageEntry) {
    this.mainPageService.saveMainPage(entry.mainPage)
      .then(() => entry.enabled = false);
  }

  public isAdmin() {
    return this.userService.isAdmin();
  }

}

export const MainPageComponent = {
  controller: MainPage,
  templateUrl: '/PageComponents/MainPage/MainPage.html',
};

export class MainPageEntry {

  public mainPage: MainPageDTO;
  public enabled = false;

  constructor(mainPage: MainPageDTO = new MainPageDTO()) {
    this.mainPage = mainPage;
  }
}
