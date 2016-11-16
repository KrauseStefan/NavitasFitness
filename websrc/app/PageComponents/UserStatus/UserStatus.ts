import { IUserDTO, UserService } from '../UserService';
import * as moment from 'moment';

const _moment: (
  inp?: moment.MomentInput, format?: moment.MomentFormatSpecification,
  language?: string, strict?: boolean
  ) => moment.Moment = (<any> moment).default;

export const statusRouterState = {
  template: '<user-status></user-status>',
  url: '/status',
};

class UserStatus {

  public model: IUserStatusModel;

  private dateFormat = 'MM/DD/YYYY';

  constructor(
    private userService: UserService,
    private $http: ng.IHttpService,
    private $location: ng.ILocationService,
    private $state: ng.ui.IStateService,
    private $q: ng.IQService) {

    this.model = {
      extendWithAmountKr: 200,
      statusMsg: 'No Active Subscription',
      transactionHistory: [],
      userEmail: '',
      validUntill: '-',
    };
    this.moniterUserCredentials();
  }

  public isLocalhost(): boolean {
    return this.$location.host().indexOf('localhost') !== -1;
  }

  public getTransactionsUpdate() {
    this.$http.get<ITransactionEntry[]>('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = res.data.map( txn => {
        txn.paymentDateParsed = _moment(txn.paymentDate).format(this.dateFormat);
        return txn;
      });
      const validTxn = this.model.transactionHistory
        .find(txn => _moment(txn.paymentDate).diff(_moment(), 'months') <= 6 )

      if (!!validTxn) {
        this.model.validUntill = _moment(validTxn.paymentDate)
          .add(6, 'months')
          .format(this.dateFormat);
        this.model.statusMsg = 'Subscription Active';
      }
    });
  }

  private moniterUserCredentials(): void {
    this.userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      if (user) {
        this.model.userEmail = user.email;
        this.getTransactionsUpdate();
      } else if (this.$state.is(statusRouterState)) {
        this.$state.go('MainPage');
      }
    }, () => this.$state.go('MainPage'));
  }

}

interface IUserStatusModel {
  userEmail: string;
  extendWithAmountKr: number;
  statusMsg: string;
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
