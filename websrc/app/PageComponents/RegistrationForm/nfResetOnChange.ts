
import IScope = angular.IScope;
import IJQuery = angular.IAugmentedJQuery;
import IDirectiveFactory = angular.IDirectiveFactory;
import INgModelController = angular.INgModelController;

const directiveName = 'nfResetOnChange';

const directiveFactoryFn: IDirectiveFactory = () => {
  return {
    link: (scope: IScope, iElement: IJQuery, iAttrs: {[att: string]: string}, ngModel: INgModelController) => {
      const errorToReset = iAttrs[directiveName];
      ngModel.$validators[errorToReset] = (modelValue: string, viewValue: string) => {
        return true;
      };
    },
    require: 'ngModel',
  };
};

export const nfResetOnChange = {
  factory: directiveFactoryFn,
  name: directiveName,
};
