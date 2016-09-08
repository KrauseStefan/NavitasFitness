import { UserService } from '../UserService';

import IHttpService = angular.IHttpService;

class UserStatus {

  public model: IUserStatusModel;

  constructor(private userService: UserService, private $http: IHttpService) {
    this.model = {
      extendWithAmountKr: 200,
      statusMsg: 'test status msg',
      transactionHistory: [],
      validUntill: '19/05-2016',
    };
    this.getTransactionsUpdate();
  }

  public getUserEmail(): string {
    if (this.userService.getLoggedinUser()) {
      return this.userService.getLoggedinUser().email;
    } else {
      return "";
    }
  }

  public getTransactionsUpdate() {
    this.$http.get('/rest/user/transactions').then((res) => {
      this.model.transactionHistory = <ITransactionEntry[]> res.data;
    });
  }
}

interface IUserStatusModel {
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
