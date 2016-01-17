import { enableProdMode } from "angular2/core"
import { bootstrap } from "angular2/platform/browser"
import { HTTP_PROVIDERS } from "angular2/http"
import { ROUTER_PROVIDERS } from "angular2/router"
import { BlogPostsService } from "./PageComponents/Blog/BlogPostsService"

import {AppComponent} from "./AppComponent"

// enableProdMode();

bootstrap(AppComponent, [
    BlogPostsService,
    HTTP_PROVIDERS,
    ROUTER_PROVIDERS
  ])