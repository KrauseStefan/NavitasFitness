import '../../Components/CkEditor/CkEditor';

export class BlogPostsService {

  private serviceUrl = 'rest/blogEntry';

  constructor(private $http: ng.IHttpService) { }

  public getBlogEntries(): ng.IPromise<BlogEntryDTO[]> {
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
  public lastEditedBy: String;
}
