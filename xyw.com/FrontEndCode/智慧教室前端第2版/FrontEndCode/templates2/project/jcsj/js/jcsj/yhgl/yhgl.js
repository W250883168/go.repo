'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 基础数据-用户管理
 */

/*    用户管理     */
app.controller("jcsjYhglContr", ['$scope','httpService','$modal', function($scope,httpService,$modal) {
	//$state.go("app.qxgl",false);
    console.log("用户管理")
    
	/*    修改密码       */
	app.controller("modalYhglPwdContr",['$scope', 'httpService', '$modalInstance','items',function ($scope, httpService,$modalInstance,items) {
		$scope.item = items;
		$scope.cancel = function() {
			$modalInstance.dismiss('cancel');
		};
	}]);

}]);

