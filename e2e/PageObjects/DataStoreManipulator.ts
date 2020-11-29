import { Datastore } from '@google-cloud/datastore';
import { RunQueryInfo, RunQueryResponse } from '@google-cloud/datastore/build/src/query';
import { get, IncomingMessage } from 'http';

export class DataManipulatorError extends Error {
  constructor(message) {
    super(message);
    this.name = 'DataManipulatorError';
  }
}

export class DataStoreManipulator {

  public static async performEmailVerification(email: string): Promise<void> {
    const [user] = await this.datastoreGetSingle('User', { 'Email': email });
    const key = user[Datastore.KEY];
    const code = await this.datastore.keyToLegacyUrlSafe(key);
    await this.sendValidationRequestFromKey(code);
  }

  public static async sendValidationRequestFromKey(code: string): Promise<void> {
    const url = 'http://localhost:8080/rest/user/verify?code=' + code;
    const response = await this.httpGet(url);

    const location = response.headers.location;
    if (location && location.includes('Verified=true')) {
      return;
    }

    throw Error('User email could not be verified using URL: ' + url);
  }

  public static async getUserId(email: string): Promise<string> {
    const [user] = await this.datastoreGetSingle('User', { 'Email': email });
    const key = user[Datastore.KEY];
    return this.datastore.keyToLegacyUrlSafe(key);
  }

  public static async getUserEntityResetSecretFromEmail(email: string): Promise<string> {
    const [user] = await this.datastoreGetSingle('User', { 'Email': email });
    if (user.PasswordResetSecret) {
      return user.PasswordResetSecret;
    }

    throw new Error(`Unable to lookup reset secret, email used: ${email}`);
  }

  public static async removeUserByAccessId(accessId: string): Promise<void> {
    const [user] = await this.datastoreQuerySingle('User', { 'AccessId': accessId });

    if (user) {
      await this.datastore.delete([user[Datastore.KEY]]);
    }
  }

  public static async removeUserByEmail(email: string): Promise<void> {
    const [user] = await this.datastoreQuerySingle('User', { 'Email': email });

    if (user) {
      await this.datastore.delete([user[Datastore.KEY]]);
    }
  }

  public static async makeUserAdmin(email: string): Promise<void> {
    const [user] = await this.datastoreQuerySingle('User', { 'Email': email });
    user.IsAdmin = true;
    this.datastore.update(user);
  }

  private static readonly datastore = new Datastore();

  private static async datastoreQuery(
    kind: string,
    filterObj: { [key: string]: string } = {},
  ): Promise<RunQueryResponse> {
    const baseQuery = this.datastore.createQuery(kind);

    const query = Object.entries(filterObj).reduce((prev, cur) => {
      const [key, value] = cur;
      return prev.filter(key, value);
    }, baseQuery);

    return this.datastore.runQuery(query);
  }

  private static async datastoreQuerySingle(
    kind: string,
    filterObj: { [key: string]: string } = {},
  ): Promise<[any | null, RunQueryInfo]> {
    const [entities, runQueryInfo] = await this.datastoreQuery(kind, filterObj);
    if (entities.length > 1) {
      throw new DataManipulatorError(`More then one '${kind}' entry was found, while only one was requested.`);
    }
    if (entities.length > 0) {
      return [entities[0], runQueryInfo];
    }
    return [null, runQueryInfo];
  }

  private static async datastoreGetSingle(
    kind: string,
    filterObj: { [key: string]: string } = {},
  ): Promise<[any, RunQueryInfo]> {
    const [entry, runQueryInfo] = await this.datastoreQuerySingle(kind, filterObj);

    if (entry === null) {
      throw new DataManipulatorError(`No '${kind}' entry was found.`);
    }
    return [entry, runQueryInfo];
  }

  private static httpGet(url: string): Promise<IncomingMessage> {
    return new Promise((resolve, reject) => {
      get(url, (res) => {
        if (res.statusCode >= 200 && res.statusCode < 400) {
          resolve(res);
        } else {
          reject(res);
        }
      });
    });
  }

}
