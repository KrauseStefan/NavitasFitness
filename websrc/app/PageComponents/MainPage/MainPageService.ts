import '../../Components/CkEditor/CkEditor';
import { MainPageDTO } from './MainPageDTO';

export class MainPageService {

  private serviceUrl = 'rest/mainPage';

  constructor(private $http: ng.IHttpService) { }

  public getMainPage(): ng.IPromise<MainPageDTO> {
    return this.$http
      .get(this.serviceUrl)
      .then((res: any) => {
        return (<MainPageDTO> res.data);
      });
  }

  public saveMainPage(mainPage: MainPageDTO) {
    const data = JSON.stringify(mainPage);

    return this.$http.put(this.serviceUrl, data);
  }

}
