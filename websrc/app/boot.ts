/// <reference path=".../../../typings/angularjs/angular.d.ts"/>
/// <reference path=".../../../typings/angular-material/angular-material"/>

angular.module( 'NavitasFitness', [ 'ngMaterial', 'ui.router' ] )
  .config(function($mdThemingProvider: angular.material.IThemingProvider) {
    $mdThemingProvider.theme('default')
      .primaryPalette('blue')
      .accentPalette('orange');
  });
  
import './AppComponent';
import './PageComponents/UserService'
