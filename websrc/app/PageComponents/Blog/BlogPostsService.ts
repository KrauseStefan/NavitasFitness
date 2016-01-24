/// <reference path=".../../../../../typings/angularjs/angular.d.ts"/>


import "../../Components/CkEditor/CkEditor"

export class BlogPostsService {

  private serviceUrl = 'rest/blogEntry'

  constructor(private $http: angular.IHttpService) { }

  public getBlogEntries(): angular.IPromise<BlogEntryDTO[]> {
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
    return this.$http.delete(this.serviceUrl + `?id=${blogEntry.Id}`)
  }
}

angular.module('NavitasFitness').service('blogPostsService', BlogPostsService)

export class BlogEntryDTO {
  Author: String
  Content: String
  Date: String
  Id: String
}