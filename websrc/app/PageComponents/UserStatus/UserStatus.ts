import { IUserDTO, UserService } from '../UserService';
import * as moment from 'moment';

export const statusRouterState: angular.ui.IState = {
  template: '<user-status></user-status>',
  url: '/status/',
};

class UserStatus {

  public model: IUserStatusModel;
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
    private $state: ng.ui.IStateService,
    private $q: ng.IQService) {

    this.model = {
      userNameStr: '',
      statusMsgKey: 'inActive',
      transactionHistory: [],
      userEmail: '',
      validUntill: '-',
    };
    this.moniterUserCredentials();
  }

  public isLocalhost(): Boolean {
    return this.$location.host().indexOf('localhost') !== -1;
  }

  public getSubmitUrl() {
    if (this.isLocalhost()) {
      return 'http://localhost:8081/processPayment';
    } else {
      return 'https://www.sandbox.paypal.com/cgi-bin/webscr';
    }
  }

  public onSubmit($event: Event) {
    const storedTimeStamp = parseInt(localStorage.getItem(this.paymentStorageKey), 10);

    if (!isNaN(storedTimeStamp) && storedTimeStamp > new Date().getTime()) {
      let confirmDialog = this.$mdDialog
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
    this.$http.get<ITransactionEntry[]>('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = res.data.map(txn => {
        txn.paymentDateParsed = moment(txn.paymentDate).format(this.dateFormat);
        return txn;
      });
      const validTxn = this.model.transactionHistory
        .filter((i) => moment(i.paymentDate).add({ month: 6 }).isSameOrAfter(moment()))
        .sort((a, b) => moment(a.paymentDate).unix())[0];

      if (!!validTxn) {
        this.model.validUntill = moment(validTxn.paymentDate)
          .add(6, 'months')
          .format(this.dateFormat);
        this.model.statusMsgKey = 'active';
      }
    });
  }

  private moniterUserCredentials(): void {
    this.userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
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

interface IUserStatusModel {
  userNameStr: string;
  userEmail: string;
  statusMsgKey: string;
  transactionHistory: Array<ITransactionEntry>;
  validUntill: string;
}

interface ITransactionEntry {
  amount: number;
  currency: string;
  paymentDate: string;
  paymentDateParsed: string;
  status: string;
}

export const UserStatusComponent = {
  controller: UserStatus,
  templateUrl: '/PageComponents/UserStatus/UserStatus.html',
};
