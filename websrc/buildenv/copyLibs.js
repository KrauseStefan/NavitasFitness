const fse = require('fs-extra');

const nodeModulesPath = './node_modules';
const appLibPath = '../webapp';

Object.keys(require('../package.json').dependencies)
    .filter(module => !(module.startsWith('@types') || module === 'material-design-icons'))
    .forEach(module => {
        const source = fse.realpathSync(`${nodeModulesPath}/${module}`);
        const target = fse.realpathSync(`${appLibPath}`) + `/libs/${module}`;
        console.log(source, '->', target);
        fse.ensureDirSync(`${target}`);
        fse.copy(source, target, {});
    });



