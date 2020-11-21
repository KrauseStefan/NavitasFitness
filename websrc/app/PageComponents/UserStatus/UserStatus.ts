import * as moment from 'moment';
import { UserService } from '../UserService';

export const statusRouterState: angular.ui.IState = {
  template: '<user-status></user-status>',
  url: '/status/',
};

class UserStatus {

  public model: UserStatusModel;
  public statusMessages: { [key: string]: string } = {
    active: 'Subscription Active',
    inActive: 'No Active Subscription',
  };

  private paymentStorageKey = 'LastPaymentClick';
  private fiveMin = 5 * 60 * 1000;

  private dateFormat = 'DD.MM.YYYY';

  constructor(
    private userService: UserService,
    private $http: ng.IHttpService,
    private $location: ng.ILocationService,
    private $mdDialog: ng.material.IDialogService,
    private $state: ng.ui.IStateService) {

    this.model = {
      userNameStr: '',
      statusMsgKey: 'inActive',
      transactionHistory: [],
      userEmail: '',
      validUntill: '-',
    };
    this.moniterUserCredentials();
  }

  public isLocalhost(): boolean {
    return this.$location.host().indexOf('localhost') !== -1;
  }

  public getSubmitUrl() {
    if (this.isLocalhost()) {
      return 'http://localhost:8082/processPayment';
    } else {
      return 'https://www.sandbox.paypal.com/cgi-bin/webscr';
    }
  }

  public onSubmit($event: Event) {
    const storedTimeStamp = this.getPaymentStorage();

    if (storedTimeStamp !== null && storedTimeStamp > new Date().getTime()) {
      const confirmDialog = this.$mdDialog
        .confirm()
        .title('Possible of multiple payments')
        .textContent(`A payment has been started recently from this browser.
        Please note that payments can take several minutes to be confirmed by paypal.
        Are you sure you want to continue?`)
        .ok('continue')
        .cancel('cancel');

      this.$mdDialog.show(confirmDialog).then(() => {
        localStorage.setItem(this.paymentStorageKey, (new Date().getTime() + this.fiveMin).toString(10));
        (<HTMLFormElement>$event.target).submit();
      });

      $event.preventDefault();
      return false;
    }

    localStorage.setItem(this.paymentStorageKey, (new Date().getTime() + this.fiveMin).toString(10));
    return true;
  }

  public getTransactionsUpdate() {
    this.$http.get<TransactionEntry[]>('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = res.data.map((txn) => {
        txn.paymentDateParsed = moment(txn.paymentDate).format(this.dateFormat);
        return txn;
      });
      const validTxn = this.model.transactionHistory
        .filter((i) => moment(i.paymentDate).add({ month: 6 }).isSameOrAfter(moment()))
        .sort((a) => moment(a.paymentDate).unix())[0];

      if (!!validTxn) {
        this.model.validUntill = moment(validTxn.paymentDate)
          .add(6, 'months')
          .format(this.dateFormat);
        this.model.statusMsgKey = 'active';
      }
    });
  }

  private getPaymentStorage() {
    const paymentStorage = localStorage.getItem(this.paymentStorageKey);
    if (paymentStorage) {
      return parseInt(paymentStorage, 10);
    }
    return null;
  }

  private moniterUserCredentials(): void {
    this.userService.getLoggedinUser$().subscribe((user) => {
      if (user) {
        this.model.userNameStr = ' - ' + user.name + ' - ' + user.accessId;
        this.model.userEmail = user.email;
        this.getTransactionsUpdate();
      } else if (this.$state.is(statusRouterState)) {
        this.$state.go('MainPage');
      }
    }, () => this.$state.go('MainPage'));
  }

}

interface UserStatusModel {
  userNameStr: string;
  userEmail: string;
  statusMsgKey: string;
  transactionHistory: TransactionEntry[];
  validUntill: string;
}

interface TransactionEntry {
  amount: number;
  currency: string;
  paymentDate: string;
  paymentDateParsed: string;
  status: string;
}

export const userStatusComponent = {
  controller: UserStatus,
  templateUrl: '/PageComponents/UserStatus/UserStatus.html',
};
