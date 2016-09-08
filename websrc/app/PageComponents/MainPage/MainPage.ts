import {BlogEntryDTO, BlogPostsService} from '../Blog/BlogPostsService';

export class MainPage {

  public entry: BlogEntry = new BlogEntry({
    author: '',
    content: '',
    date: '',
    key: null,
  });

  constructor(public blogPostsService: BlogPostsService) {
    blogPostsService.getBlogEntries()
      .then(blogEntries => this.entry = new BlogEntry(blogEntries[0]));
  }

  public saveEntry(entry: BlogEntry) {
    this.blogPostsService.saveBlogEntry(entry.blogEntry)
      .then(() => entry.enabled = false);
  }
}

export const MainPageComponent = {
  controller: MainPage,
  templateUrl: '/PageComponents/MainPage/MainPage.html',
};

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry;
  }
}
