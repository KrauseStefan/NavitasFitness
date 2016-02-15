
angular.module( 'NavitasFitness', [ 'ngMaterial', 'ui.router', 'ngCookies' ] )
  .config(function($mdThemingProvider: angular.material.IThemingProvider) {
    $mdThemingProvider.theme('default')
      .primaryPalette('blue')
      .accentPalette('orange');
  });
  
import './AppComponent';
import './PageComponents/UserService'
