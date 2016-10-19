
import IScope = angular.IScope;
import IJQuery = angular.IAugmentedJQuery;
import IDirectiveFactory = angular.IDirectiveFactory;
import INgModelController = angular.INgModelController;

const directiveName = 'nfShouldEqual';

const directiveFactoryFn: IDirectiveFactory = ($parse: angular.IParseService) => {
  return {
    link: (scope: IScope, iElement: IJQuery, attr: {[att: string]: string}, ngModel: INgModelController) => {
      const otherValue = attr[directiveName];
      const parsedExp = $parse(otherValue);
      ngModel.$validators[directiveName] = (modelValue: string, viewValue: string) => {
          return modelValue === parsedExp(scope);
      };

      scope.$watch(otherValue, () => {
          ngModel.$validate();
      });

    },
    require: 'ngModel',
  };
};

export const nfShouldEqual = {
  factory: directiveFactoryFn,
  name: directiveName,
};
