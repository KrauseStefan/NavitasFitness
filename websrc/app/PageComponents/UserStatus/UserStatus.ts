
import { UserService } from '../UserService';

import IHttpService = angular.IHttpService;

export class UserStatus {

  model: UserStatusModel;

  constructor(private userService: UserService, private $http: IHttpService) {
    this.model = {
      statusMsg: 'test status msg',
      validUntill: '19/05-2016',
      extendWithAmountKr: 200,
      transactionHistory: [],
    }
    this.getTransactionsUpdate();
  }

  public getUserEmail(): string {
    if (this.userService.getLoggedinUser()) {
      return this.userService.getLoggedinUser().email;
    } else {
      return "";
    }
  }

  getTransactionsUpdate() {
    this.$http.get('http://localhost:9000/rest/user/transactions').then(
      (res) => {
        this.model.transactionHistory = <TransactionEntry[]>res.data;
      }
    )
  }
}

class UserStatusModel {
  statusMsg: string;
  validUntill: string;
  transactionHistory: Array<TransactionEntry>;
  extendWithAmountKr: number;
}

class TransactionEntry {
  amount: number;
  currency: string;
  paymentDate: string;
  status: string;
}

angular.module('NavitasFitness')
  .component('userStatus', {
    templateUrl: '/PageComponents/UserStatus/UserStatus.html',
    controller: UserStatus
  });