import {Component, NgFor} from "angular2/angular2"
import {HTTP_PROVIDERS} from "angular2/http"
import {CkeEditorComponent} from "../ckEditor/CkEditor"
import {BlogPostsService, BlogEntry} from "./BlogPostsService"

@Component({
  selector: 'main',
  templateUrl: '/main/main.html',
  directives: [CkeEditorComponent, NgFor]
})
export class Main {

  public entries: BlogEntry[] = [];
  public enabled: boolean = false;

  constructor(blogPostsService: BlogPostsService) {

    blogPostsService.getBlogEntries()
      .subscribe(blogEntries => this.entries = (<BlogEntry[]>blogEntries.json()));

  }


}

