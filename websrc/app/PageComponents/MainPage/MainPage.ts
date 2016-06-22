import {BlogPostsService, BlogEntryDTO} from '../Blog/BlogPostsService';

export class MainPage {

  public entry: BlogEntry = new BlogEntry({
    author: '',
    content: '',
    date: '',
    key: null
  });

  constructor(public blogPostsService: BlogPostsService) {
    blogPostsService.getBlogEntries()
      .then(blogEntries => this.entry = new BlogEntry(blogEntries[0]));
  }

  saveEntry(entry: BlogEntry) {
    this.blogPostsService.saveBlogEntry(entry.blogEntry)
      .then(() => entry.enabled = false);
  }
}

export const MainPageComponent = {
  templateUrl: '/PageComponents/MainPage/MainPage.html',
  controller: MainPage
};

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry;
  }
}
