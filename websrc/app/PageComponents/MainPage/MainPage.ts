import {IUserDTO, UserService} from '../UserService';
import {MainPageDTO, MainPageService} from './MainPageService';

export class MainPage {

  public entry: MainPageEntry = new MainPageEntry({
    content: '',
    date: '',
    key: null,
    lastEditedBy: '',
  });

  public isAdmin = false;

  constructor(
    public mainPageService: MainPageService,
    private userService: UserService) {

    userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      this.isAdmin = user && user.isAdmin;
    });

    mainPageService.getMainPage()
      .then(mainPage => this.entry = new MainPageEntry(mainPage));
  }

  public saveEntry(entry: MainPageEntry) {
    this.mainPageService.saveMainPage(entry.mainPage)
      .then(() => entry.enabled = false);
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
