import { UserService } from '../UserService';
import { MainPageEntry } from './MainPageEntry';
import { MainPageService } from './MainPageService';

export const mainPageRouterState: angular.ui.IState = {
  template: '<main-page></main-page>',
  url: '/main-page/',
};

export class MainPage {

  public entry: MainPageEntry = new MainPageEntry({
    content: '',
    date: '',
    lastEditedBy: '',
  });

  public isAdmin = false;

  constructor(
    public mainPageService: MainPageService,
    userService: UserService) {

    userService.getLoggedinUser$().subscribe((user) => {
      this.isAdmin = (user && user.isAdmin) || false;
    });

    mainPageService.getMainPage()
      .then((mainPage) => this.entry = new MainPageEntry(mainPage));
  }

  public saveEntry(entry: MainPageEntry) {
    this.mainPageService.saveMainPage(entry.mainPage)
      .then(() => entry.enabled = false);
  }
}

export const mainPageComponent = {
  controller: MainPage,
  templateUrl: '/PageComponents/MainPage/MainPage.html',
};
