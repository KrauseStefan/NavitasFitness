import { Datastore } from '@google-cloud/datastore';
import { RunQueryResponse } from '@google-cloud/datastore/build/src/query';

const datastore = new Datastore();

async function dbGet(kind: string): Promise<RunQueryResponse> {
    const query = datastore
        .createQuery(kind)
        // .order('timestamp', {descending: true})
        .limit(10);

    return datastore.runQuery(query);
}

async function getKinds(): Promise<string[]> {
    const [entities] = await dbGet('__kind__');
    const keyEntities = entities.map((entity) => {
        const symbol = Object.getOwnPropertySymbols(entity)[0];
        return entity[symbol];
    });

    return keyEntities.map((entity) => entity.name);
}

async function main() {
    (await getKinds()).forEach((kind) => {
        console.log(kind);
    });
    console.log();

    const [settings] = await dbGet('Setting');
    settings.forEach((setting) => {
        const symbol = Object.getOwnPropertySymbols(setting)[0];
        delete setting[symbol];
        for (const [key, value] of Object.entries(setting)) {
            console.log(`${key} = ${value}`);
        }
        console.log();
    });

    const [users] = await dbGet('User');
    users.forEach((user) => {
        const symbol = Object.getOwnPropertySymbols(user)[0];
        delete user[symbol];
        delete user['PasswordSalt'];
        delete user['PasswordHash'];
        delete user['PasswordResetTime'];
        console.log(user);
    });

}

main();
