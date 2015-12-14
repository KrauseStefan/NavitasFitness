import {provide, Injectable, Observable} from "angular2/angular2"
import {Http, Response} from "angular2/http"
import "rxjs/add/operator/map" //This error can be ignored and will be fixed in a later version of typescript

@Injectable()
export class BlogPostsService {

  private serviceUrl = 'rest/blogEntry'

  constructor(private http: Http) { }

  public getBlogEntries(): Observable<BlogEntryDTO[]> {
    return this.http.get(this.serviceUrl)
      .map((res: Response) => (<BlogEntryDTO[]>res.json()));
  }

  saveBlogEntry(blogEntry: BlogEntryDTO) {
    const data = JSON.stringify(blogEntry);

    return this.http.put(this.serviceUrl, data);
  }

  deleteBlogEntry(blogEntry: BlogEntryDTO) {
    return this.http.delete(this.serviceUrl + `?id=${blogEntry.Id}`)
  }
}

provide('blogPostsService', { useClass: BlogPostsService });

export class BlogEntryDTO {
  Author: String
  Content: String
  Date: String
  Id: String
}