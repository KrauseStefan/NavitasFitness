import {Component} from "angular2/core"
import {HTTP_PROVIDERS} from "angular2/http"
import {CkEditorComponent} from "../../Components/CkEditor/CkEditor"
import {BlogPostsService, BlogEntryDTO} from "../Blog/BlogPostsService"

@Component({
  selector: 'main-page',
  templateUrl: '/PageComponents/MainPage/MainPage.html',
  directives: [CkEditorComponent]
})
export class MainPage {

  public entry: BlogEntry = new BlogEntry();

  constructor(public blogPostsService: BlogPostsService) {
    blogPostsService.getBlogEntries()
      .subscribe(blogEntries => {
        this.entry = new BlogEntry(blogEntries[0]);
      });
  }

  saveEntry(entry: BlogEntry) {
    this.blogPostsService.saveBlogEntry(entry.blogEntry)
      .subscribe(
      () => { }, //onNext
      () => { }, //onError
      () => entry.enabled = false //onCompleate
      )
  }
}

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry
  }
}
