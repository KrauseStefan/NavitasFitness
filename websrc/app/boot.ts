import { bootstrap } from "angular2/platform/browser"
import { HTTP_PROVIDERS } from "angular2/http"
import { BlogPostsService } from "./MainPage/BlogPostsService"

import {AppComponent} from "./AppComponent"

bootstrap(AppComponent, [
    BlogPostsService,
    HTTP_PROVIDERS
  ])