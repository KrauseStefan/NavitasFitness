
import { UserService } from '../UserService';

export class UserStatus {

  model: UserStatusModel;

  constructor(private userService: UserService) {
    this.model = {
      statusMsg: 'test status msg',
      validUntill: '19/05-2016',
      extendWithAmountKr: 200,
      transactionHistory: [
        { amountKr: 200 },
        { amountKr: 200 },
        { amountKr: 200 },
        { amountKr: 200 }
      ]
    }
  }

  public getUserEmail(): string {
    if(this.userService.getLoggedinUser()) {
      return this.userService.getLoggedinUser().email;
    } else {
      return "";
    }
  }

}

class UserStatusModel {
  statusMsg: string;
  validUntill: string;
  transactionHistory: Array<TransactionEntry>;
  extendWithAmountKr: number;
}

class TransactionEntry {
  amountKr: number;
}

angular.module('NavitasFitness')
  .component('userStatus', {
    templateUrl: '/PageComponents/UserStatus/UserStatus.html',
    controller: UserStatus
  });