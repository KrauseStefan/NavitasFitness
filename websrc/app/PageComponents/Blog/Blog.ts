import {CkEditorComponent} from '../../Components/CkEditor/CkEditor';
import {BlogPostsService, BlogEntryDTO} from './BlogPostsService';
import {UserService} from '../UserService';

export class Blog {

  public entries: BlogEntry[] = [];

  constructor(
    private blogPostsService: BlogPostsService,
    private userService: UserService) {
    blogPostsService.getBlogEntries()
      .then(blogEntries => {
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
      .then(() => entry.enabled = false)
  }

  deleteEntry(entry: BlogEntry) {
    this.blogPostsService.deleteBlogEntry(entry.blogEntry)
      .then(() => {
        const index = this.entries.indexOf(entry);
        this.entries.splice(index, 1);
      })
  }
  
  isAdmin() {
    return this.userService.isAdmin();
  }
}

export class BlogEntry {

  public blogEntry: BlogEntryDTO;
  public enabled = false;

  constructor(blogEntry: BlogEntryDTO = new BlogEntryDTO()) {
    this.blogEntry = blogEntry
  }
}

angular.module('NavitasFitness')
  .component('blog', {
    templateUrl: '/PageComponents/Blog/Blog.html',
    controller: Blog
  });