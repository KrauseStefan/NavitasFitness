import { IUserDTO, UserService } from '../UserService';
import * as moment from 'moment';

const momentFn: (
  inp?: moment.MomentInput, format?: moment.MomentFormatSpecification,
  language?: string, strict?: boolean
) => moment.Moment = (<any>moment).default;

export const statusRouterState = {
  template: '<user-status></user-status>',
  url: '/status/',
};

class UserStatus {

  public model: IUserStatusModel;

  public statusMessages: { [key: string]: string } = {
    active: 'Subscription Active',
    inActive: 'No Active Subscription',
  };

  private dateFormat = 'DD.MM.YYYY';

  constructor(
    private userService: UserService,
    private $http: ng.IHttpService,
    private $location: ng.ILocationService,
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

  public isLocalhost(): boolean {
    return this.$location.host().indexOf('localhost') !== -1;
  }

  public getTransactionsUpdate() {
    this.$http.get<ITransactionEntry[]>('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = res.data.map(txn => {
        txn.paymentDateParsed = momentFn(txn.paymentDate).format(this.dateFormat);
        return txn;
      });
      const validTxn = this.model.transactionHistory
        .filter((i) => momentFn(i.paymentDate).add({ month: 6 }).isSameOrAfter(momentFn()))
        .sort((a, b) => momentFn(a).unix())[0];

      if (!!validTxn) {
        this.model.validUntill = momentFn(validTxn.paymentDate)
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
