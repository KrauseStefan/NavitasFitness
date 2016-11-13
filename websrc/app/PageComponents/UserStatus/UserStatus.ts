import { IUserDTO, UserService } from '../UserService';

export const statusRouterState = {
  template: '<user-status></user-status>',
  url: '/status',
};

class UserStatus {

  public model: IUserStatusModel;

  constructor(
    private userService: UserService,
    private $http: ng.IHttpService,
    private $location: ng.ILocationService,
    private $state: ng.ui.IStateService,
    private $q: ng.IQService) {

    this.model = {
      extendWithAmountKr: 200,
      statusMsg: 'test status msg',
      transactionHistory: [],
      userEmail: '',
      validUntill: '19/05-2016',
    };
    this.getTransactionsUpdate();
    this.moniterUserCredentials();
  }

  public isLocalhost(): boolean {
    return this.$location.host().indexOf('localhost') !== -1;
  }

  public getTransactionsUpdate() {
    this.$http.get<ITransactionEntry[]>('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = res.data;
    });
  }

  private moniterUserCredentials(): void {
    this.userService.getLoggedinUser$().subscribe((user: IUserDTO) => {
      if (user) {
        this.model.userEmail = user.email;
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
  status: string;
}

export const UserStatusComponent = {
  controller: UserStatus,
  templateUrl: '/PageComponents/UserStatus/UserStatus.html',
};
