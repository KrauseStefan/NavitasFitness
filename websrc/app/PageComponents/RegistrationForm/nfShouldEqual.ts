
import IDirectiveFactory = angular.IDirectiveFactory;
import INgModelController = angular.INgModelController;

const directiveName = 'nfShouldEqual';

const directiveFactoryFn: IDirectiveFactory = () => {
  return {
    link: (_scope, _iElement, attr, controller) => {
      const ngModel = controller as INgModelController;
      const otherValue = attr[directiveName];
      const otherFormCtrl: INgModelController = (<any>ngModel).$$parentForm[otherValue];

      ngModel.$validators[directiveName] = (modelValue: string) => {
        return modelValue === otherFormCtrl.$viewValue;
      };

      otherFormCtrl.$validators[directiveName] = () => {
        ngModel.$validate();
        return true;
      };

    },
    require: 'ngModel',
  };
};

export const nfShouldEqual = {
  factory: directiveFactoryFn,
  name: directiveName,
};
