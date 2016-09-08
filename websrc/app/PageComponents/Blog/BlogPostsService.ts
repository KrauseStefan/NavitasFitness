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
        return (<BlogEntryDTO[]> res.data);
      });
  }

  public saveBlogEntry(blogEntry: BlogEntryDTO) {
    const data = JSON.stringify(blogEntry);

    return this.$http.put(this.serviceUrl, data);
  }

  public deleteBlogEntry(blogEntry: BlogEntryDTO) {
    return this.$http.delete(this.serviceUrl + `?id=${blogEntry.key}`);
  }
}

export class BlogEntryDTO {
  public author: String;
  public content: String;
  public date: String;
  public key: String;
}
