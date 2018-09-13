
const directiveName = 'nfResetOnChange';

function directiveFactoryFn(): ng.IDirective {

  return {
    link: (_scope, _iElement, iAttrs, controller) => {
      const ngModel = controller as ng.INgModelController;

      const errorsToResetStr = iAttrs[directiveName];
      let start = 0;
      let end = errorsToResetStr.length;
      if (errorsToResetStr[0] === '[') {
        start = 1;
      }
      if (errorsToResetStr[end - 1] === ']') {
        end = end - 1;
      }

      errorsToResetStr
        .substring(start, end)
        .split(',')
        .map((error: string) => error.trim())
        .forEach((error: string) => {
        ngModel.$validators[error] = () => {
          return true;
        };
      });
    },
    require: 'ngModel',
  };
}

export const nfResetOnChange = {
  factory: directiveFactoryFn,
  name: directiveName,
};
