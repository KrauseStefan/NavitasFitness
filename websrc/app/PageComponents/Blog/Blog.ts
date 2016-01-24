import {Component} from "angular2/core"
import {HTTP_PROVIDERS} from "angular2/http"
import {CkEditorComponent} from "../../Components/CkEditor/CkEditor"
import {BlogPostsService, BlogEntryDTO} from "./BlogPostsService"

@Component({
  selector: 'Blog',
  templateUrl: '/PageComponents/Blog/Blog.html',
  directives: [CkEditorComponent]
})
export class Blog {

  public entries: BlogEntry[] = [];

  constructor(public blogPostsService: BlogPostsService) {
    blogPostsService.getBlogEntries()
      .subscribe(blogEntries => {
        this.entries = blogEntries.map(blogEntry => {
          return new BlogEntry(blogEntry);
        });
      });
  }

  createBlogPost() {
    let entry = new BlogEntry();
    entry.enabled = true;
    this.entries.push(entry);
  }

  saveEntry(entry: BlogEntry) {
    this.blogPostsService.saveBlogEntry(entry.blogEntry)
      .subscribe(
      () => { }, //onNext
      () => { }, //onError
      () => entry.enabled = false //onCompleate
      )
  }

  deleteEntry(entry: BlogEntry) {
    this.blogPostsService.deleteBlogEntry(entry.blogEntry)
      .subscribe(
      () => { }, //onNext
      () => { }, //onError
      () => {
        const index = this.entries.indexOf(entry) //onCompleate
        this.entries.splice(index, 1);
      })
  }
}

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry
  }
}