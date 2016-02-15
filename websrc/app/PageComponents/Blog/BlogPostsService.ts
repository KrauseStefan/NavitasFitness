/// <reference path="../../../typings/main"/>
import '../../Components/CkEditor/CkEditor';

import IHttpService = angular.IHttpService;
import IPromise = angular.IPromise;

export class BlogPostsService {

  private serviceUrl = 'rest/blogEntry';

  constructor(private $http: IHttpService) { }

  public getBlogEntries(): IPromise<BlogEntryDTO[]> {
    return this.$http
      .get(this.serviceUrl)
      .then((res: any) => {
        return (<BlogEntryDTO[]>res.data)
      });
  }

  saveBlogEntry(blogEntry: BlogEntryDTO) {
    const data = JSON.stringify(blogEntry);

    return this.$http.put(this.serviceUrl, data);
  }

  deleteBlogEntry(blogEntry: BlogEntryDTO) {
    return this.$http.delete(this.serviceUrl + `?id=${blogEntry.key}`)
  }
}

angular.module('NavitasFitness').service('blogPostsService', BlogPostsService);

export class BlogEntryDTO {
  author: String
  content: String
  date: String
  key: String
}