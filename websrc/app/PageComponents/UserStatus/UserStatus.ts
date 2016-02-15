
export class UserStatus {
  
  model: UserStatusModel;
    
  constructor() {
    this.model = {
      statusMsg: 'test status msg',
      validUntill: '19/05-2016',
      extendWithAmountKr: 200,
      transactionHistory: [
        { amountKr: 200 },
        { amountKr: 200 },
        { amountKr: 200 },
        { amountKr: 200 }]
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