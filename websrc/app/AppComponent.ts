import {Component} from "angular2/core"
import {RouteConfig, Router, ROUTER_DIRECTIVES, Location} from "angular2/router"
import {MainPage} from "./PageComponents/MainPage/MainPage"
import {Blog} from "./PageComponents/Blog/Blog"
import {RegistrationForm} from "./PageComponents/RegistrationForm/RegistrationForm"

@Component({
  selector: 'app-component',
  templateUrl: './AppComponent.html',
  directives: [ ROUTER_DIRECTIVES ]
})
@RouteConfig([
  {path:'/',          name: 'MainPage', component: MainPage,  useAsDefault: true},
  {path:'/main-page', name: 'MainPage', component: MainPage},
  {path:'/blog',      name: 'Blog',     component: Blog},
  {aux: '/modal',     name: 'Modal',    component: RegistrationForm},
  {aux: '/',          name: 'None',     component: DummyComponent}
])
export class AppComponent {

  constructor(private router: Router, private location: Location) {
  }

}
