
import IScope = angular.IScope;
import IJQuery = angular.IAugmentedJQuery;
import IDirectiveFactory = angular.IDirectiveFactory;
import INgModelController = angular.INgModelController;

const directiveName = 'nfResetOnChange';

const directiveFactoryFn: IDirectiveFactory = () => {
  return {
    link: (scope: IScope, iElement: IJQuery, iAttrs: { [att: string]: string }, ngModel: INgModelController) => {
      const errorsToResetStr = iAttrs[directiveName];
      let start = 0;
      let end = errorsToResetStr.length;
      if (errorsToResetStr[0] === '[') {
        start = 1;
      }
      if (errorsToResetStr[end - 1] === ']') {
        end = end - 1;
      }

      const errorsToReset = errorsToResetStr
        .substring(start, end)
        .split(',')
        .map((error) => error.trim());

      errorsToReset.forEach((error) => {
        ngModel.$validators[error] = (modelValue: string, viewValue: string) => {
          return true;
        };
      });
    },
    require: 'ngModel',
  };
};

export const nfResetOnChange = {
  factory: directiveFactoryFn,
  name: directiveName,
};
