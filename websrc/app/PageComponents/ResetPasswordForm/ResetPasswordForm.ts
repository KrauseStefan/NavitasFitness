import { UserService } from '../UserService';

const resetMailSentMessage = 'Password reset instructions e-mail has been sent';

export class ResetPasswordFormController {

  constructor(
    private $scope: any,
    private userService: UserService,
    private $mdDialog: ng.material.IDialogService,
    private $mdToast: ng.material.IToastService) {

    $scope.model = {
      email: '',
    };

    $scope.submit = () => this.submit();
    $scope.cancel = () => this.cancel();
  }

  public submit() {
    this.userService.sendResetPasswordRequest(this.$scope.model.email)
      .then(() => {
        this.$mdToast.show(this.$mdToast.simple().textContent(resetMailSentMessage));
        this.$mdDialog.hide();
        this.displayCheckEmailNotice();
      }, (err) => {
        if (err.status === 404) {
          this.$scope.ResetPasswordForm['email'].$setValidity('NotFound', false);
        }
      });
  }

  public cancel() {
    this.$mdDialog.cancel();
  }

  private displayCheckEmailNotice() {
    return this.$mdDialog.show(
      this.$mdDialog.alert()
        .clickOutsideToClose(true)
        .title('Password rest e-mail sent')
        .textContent(`Please check your e-mail inbox for instructions on how to complete the password reset procedure.`)
        .ariaLabel('Confirmation e-mail sent')
        .ok('OK'));
  }

}
