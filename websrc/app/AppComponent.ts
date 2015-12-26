import {Component} from "angular2/core"
import {RouteConfig, Router, ROUTER_DIRECTIVES, Location} from "angular2/router"
import {MainPage} from "./PageComponents/MainPage/MainPage"
import {Blog} from "./PageComponents/Blog/Blog"

@Component({
  selector: 'app-component',
  templateUrl: './AppComponent.html',
  directives: [ROUTER_DIRECTIVES]
})
@RouteConfig([
  {path:'/',          name: 'MainPage', component: MainPage,  useAsDefault: true},
  {path:'/main-page', name: 'MainPage', component: MainPage},
  {path:'/blog',      name: 'Blog',     component: Blog},
  {aux: '/modal',     name: 'Modal',    component: MainPage}
])
export class AppComponent {

  constructor(private router: Router, private location: Location) {
  }

  //hack until aux routes gets fixed
  openRegistrerDialog() {
    // const base = window.location.pathname.split(/[\/()]/g).filter(i => i !== '')[0]
    const base = this.location.path().split(/[\/()]/g).filter(i => i !== '')[0]
    this.router.navigateByUrl(`/${base}(modal)`);
  }

}
