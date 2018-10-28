import {} from 'angular';

interface Transaction {
  amount: number;
  currency: string;
  isActive: boolean;
  paymentDate: string;
  status: string;
}

class AdminTransactionGridCtrl implements ng.IComponentController {

  public transactions: Transaction[] = [];

  private transactionsCache: { [key: string]: Transaction[] } = {};

  constructor(
    private $q: ng.IQService,
    private $http: ng.IHttpService,
  ) { }

  public $onChanges(onChangesObj: {parrentIds: ng.IChangesObject<string[]>}): void {
    if (onChangesObj.parrentIds) {
      this.displayTransactions(onChangesObj.parrentIds.currentValue);
    }
  }

  public getTransactions(key: string): ng.IPromise<Transaction[]> {
    this.transactions = [];
    return this.$http.get<Transaction[]>(`/rest/user/transactions/${key}`).then((res) => {
      this.transactionsCache[key] = res.data;
      return res.data;
    }, (resp: ng.IHttpResponse<string>) => {

      if (resp.status >= 400 && resp.status < 500) {
        return this.$q.resolve([]);
      }

      return this.$q.reject(resp.data);
    });
  }

  public async displayTransactions(selectedUserKeys: ReadonlyArray<string>) {
    const transactionsPromises = selectedUserKeys
      .map((key) => {
        const cacheHit = this.transactionsCache[key];
        if (cacheHit) {
          return this.$q.resolve(cacheHit);
        }
        return this.getTransactions(key);
      });

    const transactions = await this.$q.all(transactionsPromises);
    this.transactions = transactions.reduce((acc, val) => acc.concat(val), []); // flatten
  }

}

export const adminTransacrtionGridComponent: ng.IComponentOptions = {
  bindings: {
    parrentIds: '<',
  },
  controller: AdminTransactionGridCtrl,
  templateUrl: '/PageComponents/AdminPage/transactionsGrid/AdminTransactionGrid.html',
};
