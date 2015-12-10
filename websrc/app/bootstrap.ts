import { Component, bootstrap } from "angular2/angular2"
import { HTTP_PROVIDERS } from "angular2/http"
import { CkEditorComponent } from "./ckEditor/CkEditor"
import { Main } from "./main/Main"
import { BlogPostsService } from "./main/BlogPostsService"

bootstrap(Main, [BlogPostsService, HTTP_PROVIDERS])