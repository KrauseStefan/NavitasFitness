import {Component} from "angular2/core"
// import {RouteConfig} from "angular2/router"
import {MainPage} from "./MainPage/MainPage"

@Component({
  selector: 'app-component',
  templateUrl: '/AppComponent.html',
  directives: [MainPage]
})
export class AppComponent { }
// @RouteConfig([
//   {path:'/main-page',     name: 'Forside',    component: MainPage},
//   // {path:'/blog',          name: 'Blog',       component: HeroListComponent},
//   // {path:'/registration',  name: 'Registrer',  component: HeroDetailComponent}
// ])


