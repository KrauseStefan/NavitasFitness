import {provide, Injectable, Observable} from "angular2/angular2"
import {Http, Response} from "angular2/http"
// import { Observable } from "@reactivex/rxjs"

@Injectable()
export class BlogPostsService {

  private serviceUrl = 'rest/blogEntry'

  constructor(private http: Http) { }

  getBlogEntries() {
    //  public getBlogEntries() : Observable<Array<BlogEntry>> {
    return this.http.get(this.serviceUrl)
      // .map((res: Response) => res.json());
  }

  createBlogEntries(blogEntry: BlogEntry) {
    const data = JSON.stringify(blogEntry);

    return this.http.put(this.serviceUrl, data);
  }

  deleteBlogEntry(blogEntry: BlogEntry | String) {
    const id = (blogEntry instanceof BlogEntry) ? blogEntry.Id : blogEntry;

    return this.http.delete(this.serviceUrl + `?id=${id}`)
  }
}

provide('blogPostsService', { useClass: BlogPostsService });


export class BlogEntry {
  Author: String
  Content: String
  Date: String
  Id: String
}