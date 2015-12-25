import {Component} from "angular2/core"
import {RouteConfig, ROUTER_DIRECTIVES} from "angular2/router"
import {MainPage} from "./MainPage/MainPage"

@Component({
  selector: 'app-component',
  templateUrl: '/AppComponent.html',
  directives: [ROUTER_DIRECTIVES]
})
@RouteConfig([
  {path:'/',              name: 'FrontPage',    component: MainPage,  useAsDefault: true},
  {path:'/front-page',    name: 'FrontPage',    component: MainPage},
  {path:'/blog',          name: 'Blog',       component: MainPage},
  {path:'/registration',  name: 'Registrer',  component: MainPage}
])
export class AppComponent { }
