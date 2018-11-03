import * as moment from 'moment';

export interface AccessIdOverrideDto {
  accessId: string;
  startDate: moment.Moment;
}

export class AdminAccessOverrideRest {

  private readonly serviceUrl = '/rest/AccessIdOverride';

  constructor(private $http: ng.IHttpService) {
  }

  public getAllAccessIdOverrrides(): ng.IPromise<AccessIdOverrideDto[]> {
    return this.$http.get<AccessIdOverrideDto[]>(this.serviceUrl)
      .then((resp) => resp.data.map((d) => {
        d.startDate = moment((<any>d).startDate);
        return d;
      }));
  }

  public saveAccessIdOverride(accessIdOverride: AccessIdOverrideDto): ng.IPromise<void> {
    return this.$http.post<void>(this.serviceUrl, accessIdOverride)
      .then(() => <void>undefined);
  }

  public deleteAccessIdOverride(accessId: string): ng.IPromise<void> {
    return this.$http.delete(`${this.serviceUrl}/${accessId}`)
      .then(() => <void>undefined);
  }

}
