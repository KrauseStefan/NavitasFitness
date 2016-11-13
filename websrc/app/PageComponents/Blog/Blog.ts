import {IUserDTO, UserService} from '../UserService';
import {BlogEntryDTO, BlogPostsService} from './BlogPostsService';

export class Blog {

  public entries: BlogEntry[] = [];

  public isAdmin = false;

  constructor(
    private blogPostsService: BlogPostsService,
    private userService: UserService) {

    userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      this.isAdmin = user && user.isAdmin;
    });

    blogPostsService.getBlogEntries()
      .then(blogEntries => {
        this.entries = blogEntries.map(blogEntry => {
          return new BlogEntry(blogEntry);
        });
      });
  }

  public createBlogPost() {
    let entry = new BlogEntry();
    entry.enabled = true;
    this.entries.push(entry);
  }

  public saveEntry(entry: BlogEntry) {
    this.blogPostsService.saveBlogEntry(entry.blogEntry)
      .then(() => entry.enabled = false);
  }

  public deleteEntry(entry: BlogEntry) {
    this.blogPostsService.deleteBlogEntry(entry.blogEntry)
      .then(() => {
        const index = this.entries.indexOf(entry);
        this.entries.splice(index, 1);
      });
  }

}

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry;
  }
}

export const BlogComponent = {
  controller: Blog,
  templateUrl: '/PageComponents/Blog/Blog.html',
};
