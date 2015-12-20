import { bootstrap } from "angular2/platform/browser"
import { HTTP_PROVIDERS } from "angular2/http"
import { Main } from "./main/Main"
import { BlogPostsService } from "./main/BlogPostsService"

bootstrap(Main, [
  BlogPostsService,
  HTTP_PROVIDERS
  ])