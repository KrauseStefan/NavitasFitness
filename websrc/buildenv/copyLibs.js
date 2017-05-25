const fse = require('fs-extra');

const nodeModulesPath = './node_modules';
const appLibPath = '../App/webapp';

const filters = [
  'angular-mocks',
  'angular-sanitize',
  'material-design-icons', // not used
  'rxjs', // imported through typescript bundling
  'angular-ui-router',
];

const nodeDependencies = Object.keys(require('../package.json').dependencies);
const other = ['angular-ui-router/release', 'rxjs/bundles'];
nodeDependencies.concat(other)
  .filter(module => !module.startsWith('@types') && filters.indexOf(module) === -1)
  .forEach(module => {
    const source = fse.realpathSync(`${ nodeModulesPath }/${ module }`);
    const target = fse.realpathSync(`${ appLibPath }`) + `/libs/${ module }`;
    console.log(source, '->', target);
    fse.ensureDirSync(`${ target }`);
    fse.copy(source, target, {});
  });
