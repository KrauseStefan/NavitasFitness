import { element, module } from 'angular';

export const DefaultErrorHandlingModule = module('DefaultErrorHandling', ['ngMaterial'])
  .config(($httpProvider: ng.IHttpProvider) => {

    let instance: ng.material.IToastService = null;

    function get$mdToast(): ng.material.IToastService {
      if (instance === null) {
        instance = <ng.material.IToastService>element(document).injector().get('$mdToast');
      }
      return instance;
    }

    $httpProvider.interceptors.push(($q: ng.IQService) => {
      return {
        responseError: (response: ng.IHttpPromiseCallbackArg<any>) => {
          if (response.status >= 500 && response.status < 600) {
            const $mdToast = get$mdToast();
            const toast = $mdToast
              .simple()
              .hideDelay(0)
              .textContent("Oops! An error occurred :(")
              .highlightAction(true)
              .action("Dismiss");

            $mdToast.show(toast);
          }

          return $q.reject(response);
        },
      };
    });
  });
